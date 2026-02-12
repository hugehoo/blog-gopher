package banksalad

import (
	"strings"

	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type BankSalad struct {
}

func NewBankSalad() *BankSalad {
	return &BankSalad{}
}

var rssURL = "https://blog.banksalad.com/rss.xml"

func (b *BankSalad) CallApi() []Post {
	cfg := rss.Config{
		URL:  rssURL,
		Corp: company.BANKSALAD,
		Filter: func(item rss.Item) bool {
			// tech 블로그 글만 필터링
			return strings.Contains(item.Link, "/tech/")
		},
	}
	return rss.Parse(cfg)
}
