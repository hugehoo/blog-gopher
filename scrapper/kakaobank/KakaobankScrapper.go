package kakaobank

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var baseURL = "https://tech.kakaobank.com"
var pageURL = baseURL + "/page/"
var parsedURL, _ = url.Parse(baseURL)

func CallApi() []Post {
	var result []Post
	resultChan := make(chan []Post)
	var wg sync.WaitGroup
	maxPages := 7

	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := getPages(i)
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

func getPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = baseURL
	} else {
		url = pageURL + strconv.Itoa(page)
	}
	res, err := http.Get(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)
	doc.Find(".post").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".post-title").Text()
		date := selection.Find(".date")
		summary := selection.Find(".post-summary")
		href, _ := selection.Find(".post-title>a").Attr("href")
		parsedDate, _ := time.Parse("2006-01-02", strings.TrimSpace(date.Text()))
		post := Post{Title: strings.TrimSpace(title), Url: processUrl(href), Summary: summary.Text(), Date: parsedDate, Corp: company.KAKAOBANK}
		posts = append(posts, post)
	})
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
