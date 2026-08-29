## Overview

Once the server is running (see [INSTALL.md](INSTALL.md)), open your browser at:

```
http://localhost:8080
```

## Home Page – Bookmark List

The home page (`/`) displays all bookmarks. Each bookmark shows its URL, title, and associated tags. If no bookmarks exist, a message is displayed.

### Tag Filter

To view only bookmarks with a specific tag, append `?tag=<tagname>` to the URL. For example:

```
http://localhost:8080/?tag=dev
```

This will reload the list showing only bookmarks tagged with `dev`.

## Adding a Bookmark

Navigate to the add page:

```
http://localhost:8080/add
```

Fill in the form fields:

- **URL** (required) – the bookmark URL.
- **Title** (optional) – a descriptive title.
- **Tags** (optional) – comma-separated tags.

Click **Submit**. After successful creation, you are redirected back to the bookmark list.

## Example cURL Requests

Since the application returns HTML, you can fetch pages with `curl`:

```bash
# Get the bookmark list
curl http://localhost:8080/

# Get the add bookmark form
curl http://localhost:8080/add

# Filter by tag
curl http://localhost:8080/?tag=example
```

## Seed Data

On first start, if the bookmarks table is empty, several sample bookmarks are inserted automatically (see `database/seed.go`). This provides initial content for testing.