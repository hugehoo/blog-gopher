package kakaobank

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type Kakaobank struct {
}

func NewKakaobank() *Kakaobank {
	return &Kakaobank{}
}

var baseURL = "https://tech.kakaobank.com"
var pageURL = baseURL + "/page/"
var parsedURL, _ = url.Parse(baseURL)

func (k *Kakaobank) CallApi() []Post {
	var result []Post
	resultChan := make(chan []Post)
	var wg sync.WaitGroup
	maxPages := 7

	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := k.GetPages(i)
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

func (k *Kakaobank) GetPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = baseURL
	} else {
		url = pageURL + strconv.Itoa(page)
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
	doc.Find(".post").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".post-title a").Text()
		dateStr := selection.Find(".date").Text()
		summary := selection.Find(".post-summary").Text()
		href, exists := selection.Find(".post-title a").Attr("href")
		
		if !exists || title == "" {
			return
		}
		
		parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
		if err != nil {
			parsedDate = time.Now()
		}
		
		post := Post{
			Title:   strings.TrimSpace(title),
			Url:     processUrl(href),
			Summary: strings.TrimSpace(summary),
			Date:    parsedDate,
			Corp:    company.KAKAOBANK,
		}
		posts = append(posts, post)
	})
	//return []Post{}
	return posts
}

func processUrl(path string) string {
	if parsedPath, err := url.Parse(path); err != nil {
		fmt.Println("KakaoBank scrapper url parsing error", err)
		return baseURL
	} else {
		return parsedURL.ResolveReference(parsedPath).String()
	}
}

func normalizePath(path string) string {
	cleanPath := filepath.Clean(path)
	dir := filepath.Base(cleanPath)
	return "/posts/" + strings.TrimPrefix(dir, "/")
}
