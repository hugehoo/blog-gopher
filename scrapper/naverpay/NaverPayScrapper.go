package naverpay

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

type Naverpay struct {
}

func NewNaverpay() *Naverpay {
	return &Naverpay{}
}

var urls = []string{
	"https://medium.com/naverfinancial/fe/home",
	"https://medium.com/naverfinancial/be/home",
}

func (n *Naverpay) CallApi() []Post {
	return utils.CallGoroutineApi(n.getPages, urls)
}

func (n *Naverpay) getPages(pageURL string) []Post {

	var posts []Post
	res, err := http.Get(pageURL)

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".u-paddingTop30").Each(func(i int, selection *goquery.Selection) {
		find := selection.Find("a")
		href, _ := find.Attr("href")
		// 기존
		//title := find.Find(".u-textScreenReader").Text() // BE 카테고리의 마지막 포스트는 상위 포스트들과 다른 selector 가지는 이슈 존재
		title := find.Find("h3 > div.u-fontSize24").Text()
		summary := find.Find(".u-contentSansThin").Find(".u-fontSize18")
		date := selection.Find("time")
		if title != "" {
			date, _ := time.Parse("Jan 2, 2006", processYear(date)) // 문자열을 날짜로 파싱
			post := Post{Title: title, Url: href, Summary: summary.Text(), Date: date, Corp: company.NAVERPAY}
			posts = append(posts, post)
		}
	})
	return posts
}

func (n *Naverpay) GetPages(page int) []Post {
	// Naverpay uses a different approach with multiple URLs, so we'll return empty for now
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
