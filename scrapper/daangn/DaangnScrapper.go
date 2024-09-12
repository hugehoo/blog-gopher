package daangn

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

var baseURL = "https://medium.com/daangn/development/home"
var pageURL = baseURL

func CallApi() []Post {

	var result []Post

	// single-page blog
	pages := getPages()
	result = append(result, pages...)
	return result
}

func getPages() []Post {

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

		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text(), Corp: company.Daangn}
		posts = append(posts, post)
	})
	return posts
}
