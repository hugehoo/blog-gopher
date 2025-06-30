package kakaopay

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Kakaopay struct {
}

var baseURL = "https://tech.kakaopay.com"
var pageURL = baseURL + "/page/"

func (k *Kakaopay) CallApi() []Post {

	// 어케 totalPage 를 파악하지
	// page 범위를 넘어가면 404 를 뱉는다.

	resultChan := make(chan []Post)
	var wg sync.WaitGroup

	maxPages := 20
	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := k.GetPages(page)
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

func (k *Kakaopay) GetPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL + strconv.Itoa(page))

	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find("._postListItem_1cl5f_66>a").Each(func(i int, selection *goquery.Selection) {

		title := selection.Find("._postInfo_1cl5f_99>strong")
		href, _ := selection.Attr("href")
		summary := selection.Find("p")
		date := selection.Find("time")
		parsedDate, _ := time.Parse("2006. 1. 2", date.Text())
		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: parsedDate, Corp: company.KAKAOPAY}
		posts = append(posts, post)

	})
	return posts
}
