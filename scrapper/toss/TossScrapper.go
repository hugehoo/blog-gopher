package toss

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

var baseURL = "https://toss.tech/tech"
var postURL = "https://toss.tech"
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
	doc.Find(".css-clywuu>li>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		innerDiv := selection.Find(".css-1e3wa1f")
		title := innerDiv.Find(".typography--h6")
		summary := innerDiv.Find(".typography--p")
		originalDateString := innerDiv.Find(".typography--small").Text()
		if len(title.Text()) != 0 {
			split := strings.Split(originalDateString, "·")
			date, _ := time.Parse("2006년 1월 2일", strings.TrimSpace(split[0])) // 문자열을 날짜로 파싱
			post := Post{Title: title.Text(), Url: postURL + href, Summary: summary.Text(), Date: date.String(), Corp: company.TOSS}
			posts = append(posts, post)
		}
	})
	return posts
}
