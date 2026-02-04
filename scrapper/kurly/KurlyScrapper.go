package kurly

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Kurly struct {
}

func NewKurly() *Kurly {
	return &Kurly{}
}

var rssURL = "https://helloworld.kurly.com/feed"

func (k *Kurly) CallApi() []Post {
	return k.GetPages(1)
}

func (k *Kurly) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.KURLY,
	})
}
