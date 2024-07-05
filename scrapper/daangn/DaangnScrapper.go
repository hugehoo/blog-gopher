package daangn

import (
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

var baseURL string = "https://medium.com/daangn/development/home"
var pageURL string = baseURL

func Main() {

	// 어케 totalPage 를 파악하지
	// page 범위를 넘어가면 404 를 뱉는다.
	for i := 1; i < 2; i++ {
		pages := getPages(i)
		log.Println(pages)
	}

}

func getPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL)

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".u-xs-size12of12").Each(func(i int, selection *goquery.Selection) {
		find := selection.Find(".u-xs-marginBottom10>a")
		href, _ := find.Attr("href")
		title := find.Find("h3").Find(".u-letterSpacingTight")
		summary := find.Find(".u-contentSansThin").Find(".u-fontSize18")
		date := selection.Find("time")

		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text()}
		posts = append(posts, post)
	})
	return posts
}
