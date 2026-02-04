package woowa

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Woowa struct {
}

func NewWoowa() *Woowa {
	return &Woowa{}
}

var rssURL = "https://techblog.woowahan.com/feed"

func (w *Woowa) CallApi() []Post {
	return w.GetPages(1)
}

func (w *Woowa) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.WOOWA,
	})
}
