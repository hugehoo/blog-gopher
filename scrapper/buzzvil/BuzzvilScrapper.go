package buzzvil

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Buzzvil struct {
}

func NewBuzzvil() *Buzzvil {
	return &Buzzvil{}
}

var rssURL = "https://tech.buzzvil.com/feed"

func (b *Buzzvil) CallApi() []Post {
	return b.GetPages(1)
}

func (b *Buzzvil) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.BUZZVIL,
	})
}
