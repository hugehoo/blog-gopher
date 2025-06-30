package oliveyoung

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Oliveyoung struct {
}

const baseURL = "https://oliveyoung.tech/blog"
const OliveUrl = "https://oliveyoung.tech"

func (o *Oliveyoung) CallApi() []Post {
	var result []Post
	for i := 1; i < 11; i++ {
		pages := o.GetPages(i)
		result = append(result, pages...)
	}
	return result
}

func (o *Oliveyoung) GetPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = baseURL
	} else {
		url = baseURL + "/page/" + strconv.Itoa(page)
	}
	res, err := http.Get(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".PostList-module--item--95839>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		div := selection.Find(".PostList-module--content--de4e3")
		title := div.Find(".PostList-module--title--a2e55").Text()
		date := div.Find(".PostList-module--date--21238")
		summary := div.Find(".PostList-module--sub--424ed").Text()
		parsedDate, _ := time.Parse("2006-01-02", date.Text())
		post := Post{
			Title:   title,
			Url:     OliveUrl + href,
			Summary: summary,
			Date:    parsedDate,
			Corp:    company.OLIVE,
		}
		posts = append(posts, post)
	})
	return posts
}
