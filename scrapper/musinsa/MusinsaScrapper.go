package musinsa

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Musinsa struct {
}

func NewMusinsa() *Musinsa {
	return &Musinsa{}
}

var rssURL = "https://techblog.musinsa.com/feed"

func (m *Musinsa) CallApi() []Post {
	return m.GetPages(1)
}

func (m *Musinsa) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.MUSINSA,
	})
}
