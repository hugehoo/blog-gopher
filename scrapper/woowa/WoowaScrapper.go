package woowa

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
)

type Woowa struct {
}

func NewWoowa() *Woowa {
	return &Woowa{}
}

var baseURL = "https://techblog.woowahan.com"
var pageURL = baseURL + "/page/"

func (w *Woowa) CallApi() []Post {
	resultChan := make(chan []Post)
	var wg sync.WaitGroup

	maxPages := 47
	for i := 1; i <= maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := w.GetPages(page)
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

func (w *Woowa) GetPages(page int) []Post {
	var posts []Post
	var res *http.Response
	var err error

	client := &http.Client{}
	var req *http.Request

	if page > 1 {
		req, err = http.NewRequest("GET", pageURL+strconv.Itoa(page)+"/", nil)
	} else {
		req, err = http.NewRequest("GET", baseURL, nil)
	}
	CheckErr(err)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err = client.Do(req)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)

	doc.Find(".post-item").Each(func(i int, selection *goquery.Selection) {

		title := selection.Find(".post-title").Text()
		href, _ := selection.Find("a").Attr("href")
		summary := selection.Find(".post-excerpt").Text()
		parsedDate := getDate(selection)

		if title != "" && href != "" {
			post := Post{
				Title:   title,
				Summary: summary,
				Date:    parsedDate,
				Url:     href,
				Corp:    company.WOOWA,
			}
			posts = append(posts, post)
		}
	})

	log.Printf("Woowa: Total posts found: %d", len(posts))

	return posts
}

func getDate(selection *goquery.Selection) time.Time {
	var parsedDate time.Time
	dateStr := selection.Find(".post-author-date").Text()
	dateStr = strings.TrimSpace(dateStr)
	if dateStr != "" {
		parsedDate, _ = time.Parse("2006. 1. 2.", dateStr)
	}

	return parsedDate
}
