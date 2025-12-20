package banksalad

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"

	"github.com/PuerkitoBio/goquery"
)

type BankSalad struct {
}

func NewBankSalad() BankSalad {
	return BankSalad{}
}

var baseURL = "https://blog.banksalad.com"
var pageURL = baseURL + "/tech/page/"

func (b *BankSalad) CallApi() []Post {
	var result []Post

	// single-page blog
	for i := 1; i < 4; i++ {
		pages := b.GetPages(i)
		result = append(result, pages...)
	}
	return result
}

func (b *BankSalad) GetPages(page int) []Post {

	var posts []Post
	res, err := http.Get(pageURL + strconv.Itoa(page))

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

	year := time.Now().Year()
	prevDay := 0
	prevMonth := time.Now().Month()
	prevYear := year
	doc.Find(".postCardMinimalstyle__PostDetails-sc-12sv3cr-2").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("h2")
		href, _ := title.Find("a").Attr("href")
		summary := selection.Find(".postCardMinimalstyle__Excerpt-sc-12sv3cr-6")
		date := selection.Find(".postCardMinimalstyle__PostDate-sc-12sv3cr-3")
		day, month, err := parseDate(date.Text())
		log.Println("date:", day, month)
		if CheckErrNonFatal(err) != nil {
			return posts
		}
		// 현재 날짜에 따른 연도 계산
		calculatedYear := calculateYear(prevDay, prevMonth, prevYear, day, month)
		log.Printf("Post Date: %d %s %d", day, month.String(), calculatedYear)

		// Create a proper time.Time from the parsed components
		parsedDate := time.Date(calculatedYear, month, day, 0, 0, 0, 0, time.UTC)

		post := Post{
			Title:   title.Text(),
			Url:     baseURL + href,
			Summary: summary.Text(),
			Date:    parsedDate,
			Corp:    company.BANKSALAD}
		posts = append(posts, post)
	})
	var fake []Post
	return fake
}

// 문자열에서 월을 변환하는 함수
func parseMonth(monthStr string) (time.Month, error) {
	months := map[string]time.Month{
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
	}

	month, exists := months[strings.TrimSpace(monthStr)]
	if !exists {
		return 0, fmt.Errorf("invalid month: %s", monthStr)
	}

	return month, nil
}

// 문자열을 파싱하여 Day와 Month를 얻는 함수
func parseDate(dateStr string) (int, time.Month, error) {
	parts := strings.Split(dateStr, "  ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid date format: %s", dateStr)
	}

	// Day 파싱
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid day: %s", parts[0])
	}

	// Month 파싱
	month, err := parseMonth(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return day, month, nil
}

type Date struct {
	Day   int
	Month time.Month
}

func getYear(dates []Date, currentYear int) []int {
	years := make([]int, len(dates))
	years[0] = currentYear

	for i := 1; i < len(dates); i++ {
		if dates[i].Month > dates[i-1].Month || (dates[i].Month == dates[i-1].Month && dates[i].Day >= dates[i-1].Day) {
			// 이전 날짜보다 이후 날짜인 경우 같은 연도
			years[i] = years[i-1]
		} else {
			// 이전 날짜보다 빠른 날짜인 경우 이전 연도
			years[i] = years[i-1] - 1
		}
	}
	return years
}

// 주어진 날짜에 따라 연도 계산
func calculateYear(prevDay int, prevMonth time.Month, prevYear int, currentDay int, currentMonth time.Month) int {
	if currentMonth > prevMonth || (currentMonth == prevMonth && currentDay >= prevDay) {
		// 이전 날짜보다 미래 또는 같은 날짜면 같은 연도
		return prevYear
	}
	// 이전 날짜보다 이전 날짜면 연도를 하나 줄임
	return prevYear - 1
}
