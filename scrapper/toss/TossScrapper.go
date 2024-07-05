package toss

import (
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

var baseURL string = "https://toss.tech/tech"
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
	doc.Find(".css-clywuu>li>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		title := selection.Find(".css-1e3wa1f").Find(".typography--h6")
		summary := selection.Find(".css-1e3wa1f").Find(".typography--p")
		date := selection.Find(".css-1e3wa1f").Find(".typography--small")
		//title := selection.Find("._postInfo_1cl5f_99>strong")
		//href, _ := selection.Attr("href")
		//summary := selection.Find("p")
		//date := selection.Find("time")
		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: date.Text()}
		posts = append(posts, post)
	})
	return posts
}
