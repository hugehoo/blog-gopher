package rss

import (
	"encoding/xml"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
)

// RSS represents the RSS feed structure
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

// DateFormat constants for common RSS date formats
const (
	RFC1123  = time.RFC1123  // Mon, 02 Jan 2006 15:04:05 MST
	RFC1123Z = time.RFC1123Z // Mon, 02 Jan 2006 15:04:05 -0700
	RFC822   = time.RFC822   // 02 Jan 06 15:04 MST
	RFC822Z  = time.RFC822Z  // 02 Jan 06 15:04 -0700
)

// Config holds RSS parser configuration
type Config struct {
	URL        string
	Corp       company.Company
	DateFormat string          // Optional: custom date format
	Filter     func(Item) bool // Optional: filter function for items
}

// Parse fetches and parses RSS feed, returning posts
func Parse(cfg Config) []Post {
	var posts []Post

	client := &http.Client{}
	req, err := http.NewRequest("GET", cfg.URL, nil)
	if CheckErrNonFatal(err) != nil {
		return posts
	}

	// Set User-Agent to avoid 403 Forbidden from some servers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml, */*")

	res, err := client.Do(req)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	if CheckCodeNonFatal(res) != nil {
		return posts
	}
	defer res.Body.Close()

	var rss RSS
	decoder := xml.NewDecoder(res.Body)
	if err := decoder.Decode(&rss); err != nil {
		CheckErrNonFatal(err)
		return posts
	}

	for _, item := range rss.Channel.Items {
		if item.Title == "" {
			continue
		}

		// Apply filter if provided
		if cfg.Filter != nil && !cfg.Filter(item) {
			continue
		}

		date := parseDate(item.PubDate, cfg.DateFormat)

		// Use description as summary, fallback to extracting from content
		summary := strings.TrimSpace(item.Description)
		if summary != "" {
			// Clean HTML from description if present
			summary = extractSummary(summary, 300)
		} else if item.Content != "" {
			summary = extractSummary(item.Content, 300)
		}

		post := Post{
			Title:   sanitizeUTF8(strings.TrimSpace(item.Title)),
			Url:     strings.TrimSpace(item.Link),
			Summary: sanitizeUTF8(summary),
			Date:    date,
			Content: sanitizeUTF8(strings.TrimSpace(item.Content)),
			Corp:    cfg.Corp,
		}
		posts = append(posts, post)
	}

	return posts
}

// sanitizeUTF8 removes invalid UTF-8 characters from a string
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	// Replace invalid UTF-8 sequences with empty string
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue // Skip invalid byte
			}
		}
		v = append(v, r)
	}
	return string(v)
}

// extractSummary extracts plain text from HTML content and truncates to maxLen
func extractSummary(html string, maxLen int) string {
	// Remove HTML tags
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	text := tagRegex.ReplaceAllString(html, "")

	// Decode common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// Normalize whitespace
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Truncate to maxLen
	if len(text) > maxLen {
		text = text[:maxLen] + "..."
	}

	return text
}

// parseDate tries multiple date formats
func parseDate(dateStr string, customFormat string) time.Time {
	dateStr = strings.TrimSpace(dateStr)

	// Try custom format first
	if customFormat != "" {
		if date, err := time.Parse(customFormat, dateStr); err == nil {
			return date
		}
	}

	// Try common formats
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date
		}
	}

	return time.Now()
}
