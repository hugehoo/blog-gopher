package twonine

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Twonine struct {
}

func NewTwonine() *Twonine {
	return &Twonine{}
}

var rssURL = "https://medium.com/feed/29cm"

func (t *Twonine) CallApi() []Post {
	return t.GetPages(1)
}

func (t *Twonine) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.TWONINE,
	})
}
