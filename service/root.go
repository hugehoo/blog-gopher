package service

import (
	. "blog-gopher/common/types"
	"blog-gopher/repository"
	"blog-gopher/scrapper/daangn"
	"blog-gopher/scrapper/kakaopay"
	"blog-gopher/scrapper/oliveyoung"
	"blog-gopher/scrapper/toss"
	"sync"
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

func (s Service) UpdateAllPosts() {
	var result []Post
	result = CallGoroutineChannel(result)
	s.repo.InsertBlogs(result)
}

func (s *Service) SearchPosts(value string, page int, size int) []Post {
	return s.repo.SearchBlogs(value, page, size)
}

func CallGoroutineChannel(result []Post) []Post {
	scrapers := []func() []Post{
		kakaopay.CallApi,
		oliveyoung.CallApi,
		daangn.CallApi,
		toss.CallApi,
		//banksalad.CallApi,
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

	for posts := range resultChan {
		result = append(result, posts...)
	}

	return result
}
