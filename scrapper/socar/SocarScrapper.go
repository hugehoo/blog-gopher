package socar

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type Socar struct {
}

func NewSocar() *Socar {
	return &Socar{}
}

const baseURL = "https://tech.socarcorp.kr"

func (s *Socar) CallApi() []Post {
	var result []Post
	var wg sync.WaitGroup
	resultChan := make(chan []Post)

	maxPages := 13
	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := s.GetPages(i)
			if len(pages) > 0 {
				resultChan <- pages
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for pages := range resultChan {
		result = append(result, pages...)
	}
	return result
}

func (s *Socar) GetPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = baseURL + "/posts"
	} else {
		url = baseURL + "/posts/page" + strconv.Itoa(page)
	}
	res, err := http.Get(url)
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
	doc.Find(".post-preview").Each(func(i int, selection *goquery.Selection) {
		anchor := selection.Find("a")
		href, _ := anchor.Attr("href")
		title := anchor.Find(".post-title")
		summary := anchor.Find(".post-subtitle")
		date := selection.Find(".post-meta").Find(".date")
		parsedDate, _ := time.Parse("2006-01-02", date.Text())
		text := getContent(href)
		post := Post{Title: title.Text(), Url: baseURL + href, Summary: summary.Text(), Date: parsedDate, Content: text, Corp: company.SOCAR}
		posts = append(posts, post)
	})
	return posts
}

func getContent(href string) string {
	res, err := http.Get(baseURL + href)
	if CheckErrNonFatal(err) != nil {
		return ""
	}
	if CheckCodeNonFatal(res) != nil {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if CheckErrNonFatal(err) != nil {
		return ""
	}
	return doc.Find(".post-content").Text()
}
