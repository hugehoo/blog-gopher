package toss

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
)

type Toss struct {
}

func NewToss() *Toss {
	return &Toss{}
}

var rssURL = "https://toss.tech/rss.xml"

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
}

func (t *Toss) CallApi() []Post {
	return t.GetPages(1)
}

func (t *Toss) GetPages(_ int) []Post {
	var posts []Post
	res, err := http.Get(rssURL)

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

	// Filter engineering category articles
	for _, item := range rss.Channel.Items {
		if item.Title == "" {
			continue
		}

		// Only include engineering articles (URL contains /article/)
		if !strings.Contains(item.Link, "/article/") {
			continue
		}

		// Parse pubDate (RFC1123 format: Fri, 09 Jan 2026 04:43:00 GMT)
		date, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			date = time.Now()
		}

		post := Post{
			Title:   item.Title,
			Url:     item.Link,
			Summary: item.Description,
			Date:    date,
			Corp:    company.TOSS,
		}
		posts = append(posts, post)
	}

	return posts
}
