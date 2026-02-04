package kakaopay

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Kakaopay struct {
}

func NewKakaopay() *Kakaopay {
	return &Kakaopay{}
}

var rssURL = "https://tech.kakaopay.com/rss"

func (k *Kakaopay) CallApi() []Post {
	return k.GetPages(1)
}

func (k *Kakaopay) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.KAKAOPAY,
	})
}
