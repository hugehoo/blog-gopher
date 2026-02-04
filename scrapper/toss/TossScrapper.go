package toss

import (
	"strings"

	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Toss struct {
}

func NewToss() *Toss {
	return &Toss{}
}

var rssURL = "https://toss.tech/rss.xml"

func (t *Toss) CallApi() []Post {
	return t.GetPages(1)
}

func (t *Toss) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.TOSS,
		Filter: func(item rss.Item) bool {
			// Only include engineering articles (URL contains /article/)
			return strings.Contains(item.Link, "/article/")
		},
	})
}
