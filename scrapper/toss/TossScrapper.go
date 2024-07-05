package toss

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

var baseURL = "https://toss.tech/tech"
var pageURL = baseURL

func Main() []Post {

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
	doc.Find(".css-clywuu>li>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		innerDiv := selection.Find(".css-1e3wa1f")
		title := innerDiv.Find(".typography--h6")
		summary := innerDiv.Find(".typography--p")
		date := innerDiv.Find(".typography--small")
		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text(), Corp: company.Toss}
		posts = append(posts, post)
	})
	return posts
}
