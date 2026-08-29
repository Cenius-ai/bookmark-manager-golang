#!/usr/bin/env bash
set -euo pipefail

export PATH="/usr/local/go/bin:$PATH"

echo "=== Bookmark Manager — install ==="

# ── Go dependencies ─────────────────────────────────────────
echo "→ downloading Go modules..."
go mod download

echo "→ building binary..."
go build -o /tmp/goapp .

# ── Self-hosted fonts ──────────────────────────────────────
echo "→ downloading fonts..."
FONTS_DIR="static/fonts"
mkdir -p "$FONTS_DIR"

FONT_CSS_URL="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@12..96,400;12..96,500;12..96,600;12..96,700&family=Hanken+Grotesk:ital,wght@0,400;0,500;0,600;0,700;1,400&display=swap"

FONT_CSS=$(curl -sS -L \
  -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" \
  "$FONT_CSS_URL" 2>/dev/null || echo "")

if [ -n "$FONT_CSS" ]; then
  # Extract font URLs
  font_urls=$(echo "$FONT_CSS" | grep -oP 'url\(\Khttps://[^)]+')
  for font_url in $font_urls; do
    fname=$(basename "$font_url" | sed 's/[?].*//')
    if [ -n "$fname" ] && [ ! -f "$FONTS_DIR/$fname" ]; then
      curl -sS -L "$font_url" -o "$FONTS_DIR/$fname" 2>/dev/null || true
    fi
    # Rewrite URL to local path
    escaped_url=$(echo "$font_url" | sed 's/[\/&.]/\\&/g')
    FONT_CSS=$(echo "$FONT_CSS" | sed "s|${font_url}|/static/fonts/${fname}|g")
  done

  # Write local fonts.css with rewritten @font-face rules
  echo "$FONT_CSS" > "$FONTS_DIR/fonts.css"
  font_count=$(find "$FONTS_DIR" -type f \( -name '*.ttf' -o -name '*.woff2' -o -name '*.woff' \) 2>/dev/null | wc -l)
  echo "  fonts ready ($font_count font files)"
else
  echo "  [warn] could not download fonts — using system fallback"
fi

echo "=== Install complete ==="
echo "Run: /tmp/goapp"
