package uber

import (
	"fmt"
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

type Uber struct {
}

func NewUber() *Uber {
	return &Uber{}
}

var baseURL = "https://www.uber.com/en-KR/blog/seoul/engineering"
var pageURL = baseURL + "/page/"

func (u *Uber) CallApi() []Post {
	resultChan := make(chan []Post)
	var wg sync.WaitGroup

	maxPages := 42
	for i := 1; i <= maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := u.GetPages(page)
			fmt.Printf("page: %d %d \n", page, len(pages))
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

func (u *Uber) GetPages(page int) []Post {
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
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err = client.Do(req)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)

	//log.Printf("Uber: Starting to parse HTML for page %d", page)

	doc.Find("[data-baseweb=\"card\"][aria-label]").Each(func(i int, selection *goquery.Selection) {
		// Get title from aria-label
		title, exists := selection.Attr("aria-label")
		if !exists {
			return
		}

		// Get URL from href attribute
		href, exists := selection.Attr("href")
		if !exists {
			return
		}

		// Ensure absolute URL
		if strings.HasPrefix(href, "/") {
			href = "https://www.uber.com" + href
		}

		// Extract categories and date from the post content
		summary := ""
		var parsedDate time.Time

		// Find the date text (format: "July 2 / Global")
		dateText := selection.Find("p").Last().Text()
		parsedDate = getDate(dateText)

		// Find categories
		categoryText := selection.Find("div").First().Text()
		if categoryText != "" {
			summary = "Categories: " + categoryText
		}

		if title != "" && href != "" {
			post := Post{
				Title:   title,
				Summary: summary,
				Date:    parsedDate,
				Url:     href,
				Corp:    company.UBER,
			}
			posts = append(posts, post)
			//log.Printf("Uber: Added post: %s", title)
		}
	})

	//log.Printf("Uber: Total posts found: %d", len(posts))
	return posts
}

func getDate(dateText string) time.Time {
	var parsedDate time.Time
	dateStr := strings.TrimSpace(dateText)

	if strings.Contains(dateStr, " / ") {
		parts := strings.Split(dateStr, " / ")
		if len(parts) > 0 {
			datePart := strings.TrimSpace(parts[0])
			
			// Check if date already contains year (format: "December 11, 2024")
			if strings.Contains(datePart, ",") {
				// Format: "December 11, 2024 / Global"
				formats := []string{
					"January 2, 2006",
					"Jan 2, 2006",
				}
				
				for _, format := range formats {
					if parsed, err := time.Parse(format, datePart); err == nil {
						parsedDate = parsed
						fmt.Printf("Uber: Parsed date with year '%s' as %v\n", dateText, parsedDate)
						return parsedDate
					}
				}
			} else {
				// Format: "July 2 / Global" (no year, assume current year)
				currentYear := time.Now().Year()
				
				// Try different formats
				formats := []string{
					"January 2 2006",
					"Jan 2 2006",
				}

				// Try multiple years starting from current year going backwards
				for yearOffset := 0; yearOffset <= 2; yearOffset++ {
					testYear := currentYear - yearOffset
					fullDate := datePart + " " + strconv.Itoa(testYear)

					for _, format := range formats {
						if parsed, err := time.Parse(format, fullDate); err == nil {
							// For current year, check if date is reasonable (not too far in future)
							if yearOffset == 0 {
								// If more than 30 days in future, try previous year
								if parsed.After(time.Now().AddDate(0, 0, 30)) {
									continue
								}
							}
							parsedDate = parsed
							fmt.Printf("Uber: Parsed date without year '%s' as %v (year %d)\n", dateText, parsedDate, testYear)
							return parsedDate
						}
					}
				}
			}
		}
	}

	return parsedDate
}
