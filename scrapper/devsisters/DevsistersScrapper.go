package devsisters

import (
	company "blog-gopher/common/enum"
	"blog-gopher/common/rss"
	. "blog-gopher/common/types"
)

type Devsisters struct {
}

func NewDevsisters() *Devsisters {
	return &Devsisters{}
}

var rssURL = "https://tech.devsisters.com/rss.xml"

func (d *Devsisters) CallApi() []Post {
	return d.GetPages(1)
}

func (d *Devsisters) GetPages(_ int) []Post {
	return rss.Parse(rss.Config{
		URL:  rssURL,
		Corp: company.DEVSISTERS,
	})
}
