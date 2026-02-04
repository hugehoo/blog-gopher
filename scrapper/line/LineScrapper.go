package line

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Line struct {
}

func NewLine() *Line {
	return &Line{}
}

var rssURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"

func (l *Line) CallApi() []Post {
	return l.GetPages(1)
}

func (l *Line) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.LINE,
	})
}
