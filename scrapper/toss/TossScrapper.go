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

	if CheckErrNonFatal(err) != nil {
		return posts
	}
	if CheckCodeNonFatal(res) != nil {
		return posts
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	
	doc.Find(".css-132j2b5>li>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		title := selection.Find(".css-3lbpmg, .css-26h49c").First()
		summary := selection.Find(".css-1yz9oud, .css-17u72f3").First()
		originalDateString := selection.Find(".css-c7jnj2").Text()
		
		if len(title.Text()) != 0 {
			// "2025년 11월 14일 · 최진영" 또는 "2025년 5월 1일· 작성자" -> "2025년 11월 14일" 추출
			dateOnly := strings.Split(originalDateString, "·")[0]
			dateOnly = strings.TrimSpace(dateOnly)
			
			date, err := time.Parse("2006년 1월 2일", dateOnly)
			if err != nil {
				// 날짜 파싱 실패시 현재 시간 사용
				date = time.Now()
			}
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
