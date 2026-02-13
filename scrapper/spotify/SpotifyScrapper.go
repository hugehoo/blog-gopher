package spotify

import (
	"net/http"
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

		date := parseDateFromSelection(selection)
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

func parseDateFromSelection(selection *goquery.Selection) time.Time {
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

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"Jan 2, 2006",
		"January 2, 2006",
		"2 Jan 2006",
	}

	for _, v := range candidates {
		for _, f := range formats {
			if d, err := time.Parse(f, v); err == nil {
				return d
			}
		}
	}

	return time.Now()
}
