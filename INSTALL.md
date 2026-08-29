## 1. Prerequisites

- Go 1.22 or later

## 2. Get the Code

Obtain the project source code from its repository or download location.

## 3. Install Dependencies

```bash
go mod download
```

Alternatively, run the provided `install.sh` script, which executes `go mod download` and `go build`:

```bash
./install.sh
```

## 4. Environment Configuration

Set the following environment variables as needed:

- `PORT` – Listening port (default: 8080).
- `DATABASE_URL` – SQLite database file path (default: `bookmarks.db`).
- `FONTS_DIR` – Directory containing font files (default: `static/fonts`).
- `FONT_CSS` – Path to font CSS file (default: `static/fonts/fonts.css`).
- `FONT_CSS_URL` – URL path for font CSS (default: `/static/fonts/fonts.css`).

These can be exported in your shell or placed in a `.env` file if the application supports it.

## 5. Database Setup

The application uses SQLite. On first run, it will automatically create the database file (if it does not exist) and run migrations (see `models/setup.go`). Seed data (see `database/seed.go`) is inserted if the bookmarks table is empty. No manual database setup is required.

## 6. Run the Development Server

```bash
go run ./...
```

The server starts on `0.0.0.0:8080` (or the `PORT` value).

## 7. Running Tests

```bash
go test ./...
```

## 8. Linting

```bash
go vet ./...
```

## 9. Production Build

Build a binary:

```bash
go build ./...
```

This produces an executable named `bookmark-manager` in the current directory. Run it directly:

```bash
./bookmark-manager
```

## 10. Troubleshooting

- **Port already in use:** Change the port by setting the `PORT` environment variable.
- **SQLite database locked:** Avoid concurrent write operations; ensure only one instance is running.
- **Go version mismatch:** Verify that Go 1.22 or later is installed (`go version`).
- **Missing dependencies:** Re-run `go mod download`.