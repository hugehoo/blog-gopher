package toss

import (
	"net/http"
	"strings"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type Toss struct {
}

func NewToss() *Toss {
	return &Toss{}
}

var postURL = "https://toss.tech"
var pageURL = "https://toss.tech/tech"

func (t *Toss) CallApi() []Post {
	var result []Post

	// single-page blog
	pages := t.GetPages(1)
	result = append(result, pages...)
	return result
}

func (t *Toss) GetPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL)

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".css-132j2b5>li>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		innerDiv := selection.Find(".css-1e3wa1f")
		title := innerDiv.Find(".typography--h6")
		summary := innerDiv.Find(".typography--p")
		originalDateString := innerDiv.Find(".typography--small").Text()
		if len(title.Text()) != 0 {
			split := strings.Split(originalDateString, "·")
			date, _ := time.Parse("2006년 1월 2일", strings.TrimSpace(split[0])) // 문자열을 날짜로 파싱
			post := Post{Title: title.Text(), Url: checkUrl(href), Summary: summary.Text(), Date: date, Corp: company.TOSS}
			posts = append(posts, post)
		}
	})
	return posts
}

func checkUrl(href string) string {
	if strings.HasPrefix(href, "https://") {
		return href
	}
	return postURL + href
}
