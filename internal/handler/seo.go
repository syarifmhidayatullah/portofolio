package handler

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type SEOHandler struct {
	postSvc    service.PostService
	projectSvc service.ProjectService
	appURL     string
}

func NewSEOHandler(postSvc service.PostService, projectSvc service.ProjectService, appURL string) *SEOHandler {
	return &SEOHandler{postSvc: postSvc, projectSvc: projectSvc, appURL: strings.TrimRight(appURL, "/")}
}

func (h *SEOHandler) Robots(c *gin.Context) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /admin/

Sitemap: %s/sitemap.xml
`, h.appURL))
}

func (h *SEOHandler) Sitemap(c *gin.Context) {
	posts, _ := h.postSvc.GetAll(true)

	now := time.Now().Format("2006-01-02")

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// Static pages
	for _, page := range []struct{ loc, freq, priority string }{
		{h.appURL + "/", "weekly", "1.0"},
		{h.appURL + "/blog", "weekly", "0.9"},
		{h.appURL + "/projects", "monthly", "0.8"},
	} {
		sb.WriteString(fmt.Sprintf("  <url>\n    <loc>%s</loc>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n    <lastmod>%s</lastmod>\n  </url>\n",
			page.loc, page.freq, page.priority, now))
	}

	// Blog posts
	for _, p := range posts {
		lastmod := now
		if p.PublishedAt != nil {
			lastmod = p.PublishedAt.Format("2006-01-02")
		}
		sb.WriteString(fmt.Sprintf("  <url>\n    <loc>%s/blog/%s</loc>\n    <changefreq>monthly</changefreq>\n    <priority>0.7</priority>\n    <lastmod>%s</lastmod>\n  </url>\n",
			h.appURL, p.Slug, lastmod))
	}

	sb.WriteString("</urlset>")

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, sb.String())
}

func (h *SEOHandler) RSS(c *gin.Context) {
	posts, _ := h.postSvc.GetAll(true)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n<rss version=\"2.0\" xmlns:atom=\"http://www.w3.org/2005/Atom\">\n<channel>\n")
	sb.WriteString(fmt.Sprintf("  <title>Syarif Hidayatullah — Blog</title>\n"))
	sb.WriteString(fmt.Sprintf("  <link>%s/blog</link>\n", h.appURL))
	sb.WriteString("  <description>Thoughts on software engineering, architecture, and the craft of building things.</description>\n")
	sb.WriteString("  <language>en</language>\n")
	sb.WriteString(fmt.Sprintf("  <atom:link href=\"%s/feed.xml\" rel=\"self\" type=\"application/rss+xml\" />\n", h.appURL))
	sb.WriteString(fmt.Sprintf("  <lastBuildDate>%s</lastBuildDate>\n", time.Now().Format(time.RFC1123Z)))

	for _, p := range posts {
		pubDate := time.Now().Format(time.RFC1123Z)
		if p.PublishedAt != nil {
			pubDate = p.PublishedAt.Format(time.RFC1123Z)
		}
		sb.WriteString("  <item>\n")
		sb.WriteString(fmt.Sprintf("    <title>%s</title>\n", html.EscapeString(p.Title)))
		sb.WriteString(fmt.Sprintf("    <link>%s/blog/%s</link>\n", h.appURL, p.Slug))
		sb.WriteString(fmt.Sprintf("    <guid>%s/blog/%s</guid>\n", h.appURL, p.Slug))
		sb.WriteString(fmt.Sprintf("    <pubDate>%s</pubDate>\n", pubDate))
		if p.Excerpt != "" {
			sb.WriteString(fmt.Sprintf("    <description>%s</description>\n", html.EscapeString(p.Excerpt)))
		}
		sb.WriteString("  </item>\n")
	}

	sb.WriteString("</channel>\n</rss>")

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.String(http.StatusOK, sb.String())
}
