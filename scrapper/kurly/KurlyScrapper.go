package kurly

import (
	"net/http"
	"strconv"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type Kurly struct {
}

func NewKurly() *Kurly {
	return &Kurly{}
}

var baseURL = "https://helloworld.kurly.com"
var pageURL = baseURL

func (k *Kurly) CallApi() []Post {

	var result []Post

	// single-page blog
	pages := k.GetPages(1)
	result = append(result, pages...)
	return result
}

/*
* 당근 - 미디엄은 해당 연도에 발행된 글은 year 를 생략해서 표현함. 올해 이전에 발행된 글엔 year 를 붙인다.
 */
func (k *Kurly) GetPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL)

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)

	doc.Find(".post-card").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".post-title").Text()
		href, _ := selection.Find(".post-link").Attr("href")
		summary := selection.Find(".title-summary").Text()
		date := selection.Find(".post-date")
		if title != "" {
			date, _ := time.Parse("2006.01.02.", date.Text()) // 문자열을 날짜로 파싱
			post := Post{
				Title:   title,
				Url:     baseURL + href,
				Summary: summary,
				Date:    date,
				Corp:    company.KURLY}
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
