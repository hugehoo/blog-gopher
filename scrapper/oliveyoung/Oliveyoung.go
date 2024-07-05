package oliveyoung

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

var baseURL = "https://oliveyoung.tech"
var pageURL = baseURL + "/blog/page/"

func Main() []Post {
	var result []Post
	for i := 2; i < 3; i++ {
		pages := getPages(i)
		result = append(result, pages...)
	}
	return result
}

func getPages(page int) []Post {

	var posts []Post
	url := pageURL + strconv.Itoa(page)
	res, err := http.Get(url)
	log.Println(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".PostList-module--item--95839>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")

		div := selection.Find(".PostList-module--content--de4e3")
		title := div.Find(".PostList-module--title--a2e55")
		date := div.Find(".PostList-module--date--21238")

		post := Post{Title: title.Text(), Url: baseURL + href, Summary: "", Date: date.Text(), Corp: company.Oliveyoung}
		posts = append(posts, post)
	})
	return posts
}
