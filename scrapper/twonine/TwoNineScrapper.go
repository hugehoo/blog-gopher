package twonine

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"time"
)

var baseURL = "https://medium.com/29cm/backend/home"
var pageURL = baseURL

func CallApi() []Post {

	var result []Post

	// single-page blog
	pages := getPages()
	result = append(result, pages...)
	return result
}

/*
* 당근 - 미디엄은 해당 연도에 발행된 글은 year 를 생략해서 표현함. 올해 이전에 발행된 글엔 year 를 붙인다.
 */
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
		if title.Text() != "" {
			date, _ := time.Parse("Jan 2, 2006", processYear(date)) // 문자열을 날짜로 파싱
			post := Post{Title: title.Text(), Url: href, Summary: summary.Text(), Date: date.String(), Corp: company.TWONINE}
			posts = append(posts, post)
		}
	})
	return posts
}

func processYear(date *goquery.Selection) string {
	var temp string
	if len(date.Text()) < 8 {
		year := time.Now().Year()
		temp = date.Text() + ", " + strconv.Itoa(year)
	} else {
		temp = date.Text()
	}
	return temp
}