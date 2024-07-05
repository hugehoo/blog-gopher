package banksalad

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

var baseURL = "https://blog.banksalad.com"
var pageURL = baseURL + "/tech/page/"

func Main() []Post {

	var result []Post

	// single-page blog
	for i := 1; i < 4; i++ {
		pages := getPages(i)
		result = append(result, pages...)
	}
	return result
}

func getPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL + strconv.Itoa(page))

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".postCardMinimalstyle__PostDetails-sc-12sv3cr-2").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("h2")
		href, _ := title.Find("a").Attr("href")
		summary := selection.Find(".postCardMinimalstyle__Excerpt-sc-12sv3cr-6")
		date := selection.Find(".postCardMinimalstyle__PostDate-sc-12sv3cr-3")
		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text(), Corp: company.Banksalad}
		posts = append(posts, post)
	})
	return posts
}
