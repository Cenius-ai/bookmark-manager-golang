package database

import (
	"log"
	"strings"

	"bookmark-manager/models"

	"gorm.io/gorm"
)

type seedBookmark struct {
	URL   string
	Title string
	Tags  string
}

var demoBookmarks = []seedBookmark{
	{
		URL:   "https://go.dev/doc/effective_go",
		Title: "Effective Go — The Go Programming Language",
		Tags:  "golang,programming,documentation",
	},
	{
		URL:   "https://github.com/gofiber/fiber",
		Title: "gofiber/fiber: Express-inspired web framework for Go",
		Tags:  "golang,web,framework",
	},
	{
		URL:   "https://sqlite.org/wal.html",
		Title: "Write-Ahead Logging in SQLite",
		Tags:  "sqlite,database,performance",
	},
	{
		URL:   "https://developer.mozilla.org/en-US/docs/Web/HTML/Element",
		Title: "HTML elements reference — MDN Web Docs",
		Tags:  "html,reference,web",
	},
	{
		URL:   "https://css-tricks.com/snippets/css/a-guide-to-flexbox/",
		Title: "A Complete Guide to Flexbox — CSS-Tricks",
		Tags:  "css,layout,reference",
	},
	{
		URL:   "https://pkg.go.dev/net/http",
		Title: "Package http — Go standard library",
		Tags:  "golang,http,stdlib",
	},
	{
		URL:   "https://12factor.net/",
		Title: "The Twelve-Factor App — Methodology for SaaS",
		Tags:  "devops,architecture,best-practices",
	},
	{
		URL:   "https://caniuse.com/",
		Title: "Can I use... Browser support tables for web features",
		Tags:  "css,html,browser,compatibility",
	},
	{
		URL:   "https://github.com/glebarez/sqlite",
		Title: "glebarez/sqlite: Pure-Go SQLite driver for GORM",
		Tags:  "golang,sqlite,gorm",
	},
	{
		URL:   "https://restfulapi.net/",
		Title: "REST API Tutorial — Resource guide to REST APIs",
		Tags:  "api,rest,architecture",
	},
	{
		URL:   "https://html.spec.whatwg.org/multipage/",
		Title: "HTML Living Standard — WHATWG",
		Tags:  "html,spec,web",
	},
	{
		URL:   "https://daringfireball.net/projects/markdown/",
		Title: "Daring Fireball: Markdown Syntax Guide",
		Tags:  "markdown,reference,documentation",
	},
}

func Seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Bookmark{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("database already seeded, skipping")
		return nil
	}

	log.Println("seeding demo bookmarks...")

	for _, sb := range demoBookmarks {
		bm := models.Bookmark{
			URL:   sb.URL,
			Title: sb.Title,
		}

		tagNames := strings.Split(sb.Tags, ",")
		for _, tn := range tagNames {
			tn = strings.TrimSpace(strings.ToLower(tn))
			if tn == "" {
				continue
			}
			var tag models.Tag
			if err := db.Where("name = ?", tn).FirstOrCreate(&tag, models.Tag{Name: tn}).Error; err != nil {
				return err
			}
			bm.Tags = append(bm.Tags, tag)
		}

		if err := db.Create(&bm).Error; err != nil {
			return err
		}
	}

	log.Printf("seeded %d bookmarks with tags\n", len(demoBookmarks))
	return nil
}
