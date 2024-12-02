package buzzvil

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var baseURL = "https://tech.buzzvil.com"
var pageURL = baseURL + "/page/"

func CallApi() []Post {

	// 어케 totalPage 를 파악하지
	// page 범위를 넘어가면 404 를 뱉는다.

	resultChan := make(chan []Post)
	var wg sync.WaitGroup

	maxPages := 11
	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := getPages(page)
			if len(pages) > 0 {
				resultChan <- pages
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()
	var result []Post
	for pages := range resultChan {
		result = append(result, pages...)
	}
	return result
}

func getPages(page int) []Post {

	var posts []Post
	var res *http.Response
	var err error

	if page > 1 {
		res, err = http.Get(pageURL + strconv.Itoa(page))
	} else {
		res, err = http.Get(baseURL)
	}

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find("article").Each(func(i int, selection *goquery.Selection) {
		// 제목 텍스트 추출
		title := selection.Find("a.post-title")
		href, _ := selection.Find(".post-title").Attr("href")
		summary := selection.Find("p")
		parsedDate := getDate(selection)
		post := Post{
			Title:   title.Text(),
			Summary: summary.Text(),
			Date:    parsedDate,
			Url:     href,
			Corp:    company.BUZZVIL}
		posts = append(posts, post)
	})
	return posts
}

func getDate(selection *goquery.Selection) time.Time {
	var parsedDate time.Time
	val := selection.Find("ul.card-meta span")
	val.Each(func(i int, selection *goquery.Selection) {
		length := val.Length()
		if i == length-1 { // author 가 여럿인 경우가 있어 날짜의 순서가 고정되지 않기 때문에 i 값으로 하드코딩하여 분기하는건 X
			date := selection.Text()
			parsedDate, _ = time.Parse("2 Jan, 2006", date)
		}
	})
	return parsedDate
}
