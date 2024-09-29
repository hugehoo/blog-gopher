package kakaobank

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var baseURL = "https://tech.kakaobank.com"
var pageURL = baseURL + "/page/"

func CallApi() []Post {
	var result []Post
	for i := 1; i < 7; i++ {
		pages := getPages(i)
		result = append(result, pages...)
	}
	return result
}

func getPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = baseURL
	} else {
		url = pageURL + strconv.Itoa(page)
	}
	res, err := http.Get(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".post").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".post-title").Text()
		date := selection.Find(".date")
		summary := selection.Find(".post-summary")
		href, _ := selection.Find(".post-title>a").Attr("href")
		parsedDate, _ := time.Parse("2006-01-02", strings.TrimSpace(date.Text()))
		post := Post{Title: strings.TrimSpace(title), Url: baseURL + strings.TrimSpace(href), Summary: summary.Text(), Date: parsedDate.String(), Corp: company.KAKAOBANK}
		posts = append(posts, post)
	})
	return posts
}
