package daangn

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Daangn struct {
}

func NewDaangn() *Daangn {
	return &Daangn{}
}

var rssURL = "https://medium.com/feed/daangn"

func (d *Daangn) CallApi() []Post {
	return d.GetPages(1)
}

func (d *Daangn) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.DAANGN,
	})
}
