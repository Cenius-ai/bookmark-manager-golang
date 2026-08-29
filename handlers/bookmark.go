package handlers

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"bookmark-manager/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// tagNameRe allows lowercase letters, digits, hyphens, underscores, and dots.
// Must be 1–40 characters. Blocks HTML metacharacters and whitespace.
var tagNameRe = regexp.MustCompile(`^[a-z0-9._-]{1,40}$`)

type BookmarkHandler struct {
	DB *gorm.DB
}

func (h *BookmarkHandler) List(c *fiber.Ctx) error {
	tagFilter := strings.TrimSpace(c.Query("tag"))

	var bookmarks []models.Bookmark
	q := h.DB.Preload("Tags").Order("bookmarks.created_at DESC")

	if tagFilter != "" {
		q = q.Joins("JOIN bookmark_tags ON bookmark_tags.bookmark_id = bookmarks.id").
			Joins("JOIN tags ON tags.id = bookmark_tags.tag_id").
			Where("tags.name = ?", tagFilter)
	}

	if err := q.Find(&bookmarks).Error; err != nil {
		log.Printf("error loading bookmarks: %v", err)
		return c.Status(500).SendString("could not load bookmarks")
	}

	var allTags []models.Tag
	if err := h.DB.Order("name ASC").Find(&allTags).Error; err != nil {
		log.Printf("error loading tags: %v", err)
		return c.Status(500).SendString("could not load tags")
	}

	return c.Render("list", fiber.Map{
		"Bookmarks": bookmarks,
		"AllTags":   allTags,
		"TagFilter": tagFilter,
		"Title":     "All Bookmarks",
	})
}

func (h *BookmarkHandler) AddForm(c *fiber.Ctx) error {
	var allTags []models.Tag
	if err := h.DB.Order("name ASC").Find(&allTags).Error; err != nil {
		log.Printf("error loading tags: %v", err)
		return c.Status(500).SendString("could not load tags")
	}

	return c.Render("add", fiber.Map{
		"AllTags":   allTags,
		"Title":     "Add Bookmark",
		"Error":     "",
		"FormURL":   "",
		"FormTitle": "",
		"FormTags":  "",
	})
}

func (h *BookmarkHandler) Create(c *fiber.Ctx) error {
	formURL := strings.TrimSpace(c.FormValue("url"))
	formTitle := strings.TrimSpace(c.FormValue("title"))
	formTags := strings.TrimSpace(c.FormValue("tags"))

	var allTags []models.Tag
	h.DB.Order("name ASC").Find(&allTags)

	renderErr := func(msg string) error {
		return c.Render("add", fiber.Map{
			"AllTags":   allTags,
			"Title":     "Add Bookmark",
			"Error":     msg,
			"FormURL":   formURL,
			"FormTitle": formTitle,
			"FormTags":  formTags,
		})
	}

	// ── URL validation ──────────────────────────────────────
	if formURL == "" {
		return renderErr("URL is required.")
	}
	if len(formURL) > 2048 {
		return renderErr("URL is too long (max 2048 characters).")
	}
	parsed, err := url.Parse(formURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return renderErr("Please enter a valid URL starting with http:// or https://")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return renderErr("URL scheme must be http or https.")
	}

	// ── Title validation ────────────────────────────────────
	if formTitle == "" {
		return renderErr("Title is required.")
	}
	if len(formTitle) > 500 {
		return renderErr("Title is too long (max 500 characters).")
	}
	// Collapse multiple whitespace into single space
	formTitle = strings.Join(strings.Fields(formTitle), " ")

	// ── Tag sanitisation ────────────────────────────────────
	bm := models.Bookmark{
		URL:   formURL,
		Title: formTitle,
	}

	if formTags != "" {
		tagNames := strings.Split(formTags, ",")
		for _, tn := range tagNames {
			tn = strings.TrimSpace(tn)
			if tn == "" {
				continue
			}
			// Normalise to lowercase and validate charset + length
			tn = strings.ToLower(tn)
			if !tagNameRe.MatchString(tn) {
				return renderErr("Invalid tag \"" + tn + "\". Tags may only contain letters, digits, hyphens, underscores and dots (1–40 chars).")
			}

			var tag models.Tag
			if err := h.DB.Where("name = ?", tn).FirstOrCreate(&tag, models.Tag{Name: tn}).Error; err != nil {
				log.Printf("error creating tag %q: %v", tn, err)
				return renderErr("Could not save tags.")
			}
			bm.Tags = append(bm.Tags, tag)
		}
	}

	if err := h.DB.Create(&bm).Error; err != nil {
		log.Printf("error creating bookmark: %v", err)
		return renderErr("Could not save bookmark.")
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}
