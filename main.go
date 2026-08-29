package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bookmark-manager/database"
	"bookmark-manager/handlers"
	"bookmark-manager/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

func main() {
	// ── Database ────────────────────────────────────────────
	db, err := database.Open()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	// ── Template engine ─────────────────────────────────────
	templateSub, err := fs.Sub(templateFS, "templates")
	if err != nil {
		log.Fatalf("failed to create template sub-fs: %v", err)
	}

	engine := html.NewFileSystem(http.FS(templateSub), ".html")

	// ── Fiber app ───────────────────────────────────────────
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layout",
	})

	app.Use(recover.New())
	app.Use(securityHeaders)

	// ── Static assets served from embed.FS ──────────────────
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("failed to create static sub-fs: %v", err)
	}
	app.Get("/static/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		data, err := fs.ReadFile(staticSub, path)
		if err != nil {
			return c.Status(404).SendString("not found")
		}
		c.Type(path)
		return c.Send(data)
	})

	// ── Handlers & routes ───────────────────────────────────
	h := &handlers.BookmarkHandler{DB: db}

	app.Get("/", h.List)
	app.Get("/add", h.AddForm)
	app.Post("/bookmarks", h.Create)

	// ── Shutdown handling ───────────────────────────────────
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("shutting down...")
		if err := app.Shutdown(); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	// ── Start server ────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "0.0.0.0:" + port
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func securityHeaders(c *fiber.Ctx) error {
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; img-src 'self' data:;")
	return c.Next()
}
