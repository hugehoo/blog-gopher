package service

import (
	. "blog-gopher/common/types"
	"blog-gopher/repository"
	"blog-gopher/scrapper/bucketplace"
	"blog-gopher/scrapper/buzzvil"
	"blog-gopher/scrapper/daangn"
	"blog-gopher/scrapper/devsisters"
	"blog-gopher/scrapper/kakaobank"
	"blog-gopher/scrapper/kakaopay"
	"blog-gopher/scrapper/kurly"
	"blog-gopher/scrapper/line"
	"blog-gopher/scrapper/musinsa"
	"blog-gopher/scrapper/naverpay"
	"blog-gopher/scrapper/oliveyoung"
	"blog-gopher/scrapper/socar"
	"blog-gopher/scrapper/toss"
	"blog-gopher/scrapper/twonine"
	"log"
	"sync"
	"time"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repo: repository}
}

func (s *Service) GetPosts(category []string, page int, pageSize int) []Post {
	return s.repo.FindBlogs(category, page, pageSize)
}

func (s *Service) SearchPostsById(id string) Post {
	return s.repo.SearchPostById(id)
}

func (s Service) UpdateAllPosts() {
	result := CallGoroutineChannel()
	s.repo.InsertBlogs(result)
}

func (s Service) UpdateLatestPosts() {
	result := CallGoroutineChannel()

	savedLatestDate := s.repo.GetLatestPost()
	parsedLatestDate, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", savedLatestDate)
	var filterResult []Post
	for _, res := range result {
		date, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", res.Date)
		if date.After(parsedLatestDate) {
			filterResult = append(filterResult, res)
		}
	}
	s.repo.InsertBlogs(filterResult)
}

func (s *Service) SearchPosts(value string, page int, size int) []Post {
	return s.repo.SearchBlogs(value)
}

func (s *Service) UpdatePost(postID string, text string) {
	_, err := s.repo.UpdatePost(postID, text)
	if err != nil {
		log.Printf("Error updating post %s: %v", postID, err)
	} else {
		log.Printf("Successfully updated post %s", postID)
	}
}

func CallGoroutineChannel() []Post {
	scrapers := []func() []Post{
		naverpay.CallApi,
		socar.CallApi,
		kakaopay.CallApi,
		kakaobank.CallApi,
		oliveyoung.CallApi,
		daangn.CallApi,
		toss.CallApi,
		musinsa.CallApi,
		twonine.CallApi,
		buzzvil.CallApi,
		kurly.CallApi,
		devsisters.CallApi,
		bucketplace.CallApi,
		line.CallApi,
	}
	resultChan := make(chan []Post, len(scrapers))

	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(scrapeFunc func() []Post) {
			defer wg.Done()
			resultChan <- scrapeFunc()
		}(scraper)
	}

	// 결과 수집을 위한 고루틴
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var result []Post
	for posts := range resultChan {
		result = append(result, posts...)
	}

	return result
}
