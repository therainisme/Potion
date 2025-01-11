# Potion

A lightweight proxy tool for customizing Notion pages, allowing you to:

1. Access Notion pages through your own server, support for custom domain names
2. Customize page title and description

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/therainisme/potion.git
cd potion
```

2. Create and edit `.env` file:

```env
PORT=8080
SITE_DOMAIN=https://your-blog.notion.site
SITE_SLUG=your-page-slug
PAGE_TITLE=Your Custom Title
PAGE_DESCRIPTION=Your Custom Description
```

3. Run:

```bash
go run .
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| SITE_DOMAIN | Your Notion site domain | https://your-blog.notion.site |
| SITE_SLUG | Page slug | your-page-slug |
| PAGE_TITLE | Custom page title | My Blog |
| PAGE_DESCRIPTION | Custom page description | Welcome to my blog |

