package daangn

import (
	"net/http"
	"strconv"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"blog-gopher/common/utils"

	"github.com/PuerkitoBio/goquery"
)

type Daangn struct {
}

func NewDaangn() *Daangn {
	return &Daangn{}
}

var urls = []string{
	"https://medium.com/daangn/development/home",
	"https://medium.com/daangn/search-home/home",
	"https://medium.com/daangn/data/home",
	"https://medium.com/daangn/machine-learning/home",
}

func (d *Daangn) CallApi() []Post {
	return utils.CallGoroutineApi(d.getPages, urls)
}

/*
* 당근 - 미디엄은 해당 연도에 발행된 글은 year 를 생략해서 표현함. 올해 이전에 발행된 글엔 year 를 붙인다.
 */
func (d *Daangn) getPages(pageURL string) []Post {

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
			post := Post{Title: title.Text(), Url: href, Summary: summary.Text(), Date: date, Corp: company.DAANGN}
			posts = append(posts, post)
		}
	})
	return posts
}

func (d *Daangn) GetPages(page int) []Post {
	// Daangn uses a different approach with multiple URLs, so we'll return empty for now
	// The actual implementation uses getPages with specific URLs
	return []Post{}
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
