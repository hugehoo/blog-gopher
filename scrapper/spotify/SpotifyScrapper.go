package spotify

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type Spotify struct {
}

func NewSpotify() *Spotify {
	return &Spotify{}
}

var baseURL = "https://confidence.spotify.com/blog"

var rssURLs = []string{
	"https://confidence.spotify.com/blog/rss.xml",
	"https://confidence.spotify.com/blog/feed.xml",
	"https://confidence.spotify.com/rss.xml",
	"https://confidence.spotify.com/feed.xml",
}

func (s *Spotify) CallApi() []Post {
	if posts := s.getFromRSS(); len(posts) > 0 {
		return posts
	}
	return s.GetPages(1)
}

func (s *Spotify) GetPages(_ int) []Post {
	var posts []Post

	req, err := http.NewRequest("GET", baseURL, nil)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	if CheckCodeNonFatal(res) != nil {
		return posts
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if CheckErrNonFatal(err) != nil {
		return posts
	}

	// Prefer structured metadata when available to get stable title/date/url.
	if structuredPosts := parsePostsFromJSONLD(doc); len(structuredPosts) > 0 {
		return structuredPosts
	}

	seen := map[string]struct{}{}
	doc.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if !exists {
			return
		}

		url := normalizeURL(href)
		if url == "" || !strings.Contains(url, "/blog/") || url == baseURL {
			return
		}

		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}

		title := strings.TrimSpace(selection.Text())
		if title == "" {
			title = strings.TrimSpace(selection.Find("h1,h2,h3,h4").First().Text())
		}
		if title == "" {
			return
		}

		summary := strings.TrimSpace(selection.Find("p").First().Text())
		if summary == "" {
			summary = title
		}

		date, ok := parseDateFromSelection(selection)
		if !ok {
			if resolved, found := s.fetchDateFromPostURL(url); found {
				date = resolved
			} else {
				date = time.Now()
			}
		}
		posts = append(posts, Post{
			Title:   title,
			Url:     url,
			Summary: summary,
			Date:    date,
			Corp:    company.SPOTIFY,
		})
	})

	return posts
}

func (s *Spotify) getFromRSS() []Post {
	for _, rssURL := range rssURLs {
		posts := rss.Parse(rss.Config{
			URL:         rssURL,
			Corp:        company.SPOTIFY,
			Suppress404: true,
			Filter: func(item rss.Item) bool {
				return strings.Contains(item.Link, "/blog/")
			},
		})
		if len(posts) > 0 {
			return posts
		}
	}
	return nil
}

func normalizeURL(href string) string {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "#") {
		return ""
	}
	if strings.HasPrefix(href, "/") {
		return "https://confidence.spotify.com" + href
	}
	if strings.HasPrefix(href, "https://confidence.spotify.com/") {
		return href
	}
	return ""
}

func parseDateFromSelection(selection *goquery.Selection) (time.Time, bool) {
	candidates := []string{}

	if datetime, ok := selection.Attr("datetime"); ok {
		candidates = append(candidates, strings.TrimSpace(datetime))
	}

	selection.Find("time").Each(func(i int, s *goquery.Selection) {
		if datetime, ok := s.Attr("datetime"); ok {
			candidates = append(candidates, strings.TrimSpace(datetime))
		}
		if txt := strings.TrimSpace(s.Text()); txt != "" {
			candidates = append(candidates, txt)
		}
	})

	// Spotify cards can render date as plain text instead of <time>.
	candidates = append(candidates, extractDateCandidates(selection.Text())...)
	candidates = append(candidates, extractDateCandidates(selection.Parent().Text())...)
	candidates = append(candidates, extractDateCandidates(selection.Closest("article,li,section,div").Text())...)

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"Jan 2, 2006",
		"January 2, 2006",
		"2 Jan 2006",
	}

	for _, v := range candidates {
		if parsed, ok := tryParseDate(v, formats); ok {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func parsePostsFromJSONLD(doc *goquery.Document) []Post {
	var posts []Post
	seen := map[string]struct{}{}

	doc.Find(`script[type="application/ld+json"]`).Each(func(i int, s *goquery.Selection) {
		payload := strings.TrimSpace(s.Text())
		if payload == "" {
			return
		}

		var root interface{}
		if err := json.Unmarshal([]byte(payload), &root); err != nil {
			return
		}

		nodes := flattenJSONLDNodes(root)
		for _, node := range nodes {
			typ := strings.ToLower(readString(node, "@type"))
			if typ == "" || (!strings.Contains(typ, "blogposting") && !strings.Contains(typ, "article")) {
				continue
			}

			url := normalizeURL(readString(node, "url"))
			if url == "" || !strings.Contains(url, "/blog/") {
				continue
			}

			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}

			title := strings.TrimSpace(readString(node, "headline"))
			if title == "" {
				title = strings.TrimSpace(readString(node, "name"))
			}
			if title == "" {
				continue
			}

			summary := strings.TrimSpace(readString(node, "description"))
			if summary == "" {
				summary = title
			}

			dateRaw := strings.TrimSpace(readString(node, "datePublished"))
			date := time.Now()
			if parsed, ok := tryParseDate(dateRaw, []string{
				time.RFC3339,
				time.RFC3339Nano,
				"2006-01-02",
				"Jan 2, 2006",
				"January 2, 2006",
				"2 Jan 2006",
			}); ok {
				date = parsed
			}

			posts = append(posts, Post{
				Title:   title,
				Url:     url,
				Summary: summary,
				Date:    date,
				Corp:    company.SPOTIFY,
			})
		}
	})

	return posts
}

func flattenJSONLDNodes(v interface{}) []map[string]interface{} {
	var out []map[string]interface{}

	switch t := v.(type) {
	case map[string]interface{}:
		out = append(out, t)
		for _, key := range []string{"@graph", "itemListElement", "mainEntity"} {
			if nested, ok := t[key]; ok {
				out = append(out, flattenJSONLDNodes(nested)...)
			}
		}
	case []interface{}:
		for _, item := range t {
			out = append(out, flattenJSONLDNodes(item)...)
		}
	}

	return out
}

func readString(node map[string]interface{}, key string) string {
	raw, ok := node[key]
	if !ok || raw == nil {
		return ""
	}

	switch v := raw.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Some fields can be nested objects with @id.
		if id, ok := v["@id"].(string); ok {
			return id
		}
	}

	return ""
}

func tryParseDate(v string, formats []string) (time.Time, bool) {
	if v == "" {
		return time.Time{}, false
	}

	v = strings.TrimSpace(v)
	for _, f := range formats {
		if d, err := time.Parse(f, v); err == nil {
			return d, true
		}
	}
	return time.Time{}, false
}

func (s *Spotify) fetchDateFromPostURL(postURL string) (time.Time, bool) {
	req, err := http.NewRequest("GET", postURL, nil)
	if err != nil {
		return time.Time{}, false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return time.Time{}, false
	}
	if res.StatusCode != http.StatusOK {
		return time.Time{}, false
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return time.Time{}, false
	}

	// 1) Try JSON-LD on article page.
	if posts := parsePostsFromJSONLD(doc); len(posts) > 0 {
		for _, p := range posts {
			if p.Url == postURL && !p.Date.IsZero() {
				return p.Date, true
			}
		}
		if !posts[0].Date.IsZero() {
			return posts[0].Date, true
		}
	}

	// 2) Common metadata tags.
	metaSelectors := []string{
		`meta[property="article:published_time"]`,
		`meta[name="article:published_time"]`,
		`meta[property="og:published_time"]`,
		`meta[name="pubdate"]`,
		`meta[name="publish_date"]`,
	}
	for _, sel := range metaSelectors {
		if content, ok := doc.Find(sel).Attr("content"); ok {
			if d, parsed := tryParseDate(content, []string{
				time.RFC3339,
				time.RFC3339Nano,
				"2006-01-02",
				"Jan 2, 2006",
				"January 2, 2006",
				"2 Jan 2006",
			}); parsed {
				return d, true
			}
		}
	}

	// 3) Visible <time> content.
	if d, ok := parseDateFromSelection(doc.Find("time").First()); ok {
		return d, true
	}

	return time.Time{}, false
}

func extractDateCandidates(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	patterns := []string{
		`\b[A-Z][a-z]{2} \d{1,2}, \d{4}\b`, // Jan 26, 2026
		`\b[A-Z][a-z]+ \d{1,2}, \d{4}\b`,   // January 26, 2026
		`\b\d{4}-\d{2}-\d{2}\b`,            // 2026-01-26
		`\b[A-Z][a-z]{2} \d{1,2} \d{4}\b`,  // Jan 26 2026
		`\b[A-Z][a-z]+ \d{1,2} \d{4}\b`,    // January 26 2026
	}

	var result []string
	seen := map[string]struct{}{}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(text, -1)
		for _, m := range matches {
			if _, ok := seen[m]; ok {
				continue
			}
			seen[m] = struct{}{}
			result = append(result, m)
		}
	}

	return result
}
