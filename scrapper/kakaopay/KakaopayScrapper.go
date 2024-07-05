package kakaopay

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

var baseURL = "https://tech.kakaopay.com"
var pageURL = baseURL + "/page/"

func Main() []Post {

	// 어케 totalPage 를 파악하지
	// page 범위를 넘어가면 404 를 뱉는다.
	var result []Post
	for i := 1; i < 2; i++ {
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
	doc.Find("._postListItem_1cl5f_66>a").Each(func(i int, selection *goquery.Selection) {

		title := selection.Find("._postInfo_1cl5f_99>strong")
		href, _ := selection.Attr("href")
		summary := selection.Find("p")
		date := selection.Find("time")

		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text(), Corp: company.Kakaopay}
		posts = append(posts, post)

	})
	return posts
}
