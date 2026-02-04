package socar

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Socar struct {
}

func NewSocar() *Socar {
	return &Socar{}
}

var rssURL = "https://tech.socarcorp.kr/feed.xml"

func (s *Socar) CallApi() []Post {
	return s.GetPages(1)
}

func (s *Socar) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.SOCAR,
	})
}
