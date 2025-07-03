package service

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"

	"blog-gopher/common/dto"
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
	"blog-gopher/scrapper/woowa"
)

type Service struct {
	repo *repository.Repository
}

type Category struct {
	Category string `form:"category"`
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repo: repository}
}

func (s *Service) GetPosts(c *gin.Context) []dto.PostDTO {
	start := time.Now()

	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	// 기본 값 설정
	page := 1
	pageSize := 40

	// 쿼리 파라미터에서 페이지 번호와 페이지 크기를 정수로 변환
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = ps
		}
	}

	var category Category
	if err := c.ShouldBind(&category); err != nil {
		log.Fatalln("Something wrong", err)
	}
	categoryStrings := category.Category
	var splits []string
	if categoryStrings != "" {
		splits = strings.Split(categoryStrings, ",")
	} else {
		splits = []string{}
	}
	blogs := s.repo.FindBlogs(splits, page, pageSize)
	response := make([]dto.PostDTO, len(blogs))
	for idx, result := range blogs {
		response[idx] = dto.ConvertToDTO(result)
	}

	end := time.Since(start)
	log.Println("total Exe:", end)
	return response
}

func (s *Service) SearchPostsById(id string) Post {
	return s.repo.SearchPostById(id)
}

func (s *Service) AutoSearchPosts(keyword string) []dto.PostDTO {
	queryResult := s.repo.AutoSearchQuery(keyword)
	var response []dto.PostDTO
	for _, result := range queryResult {
		response = append(response, dto.ConvertToDTO(result))
	}
	return response
}

func (s Service) UpdateAllPosts() {
	result := CallGoroutineChannel()
	s.repo.InsertBlogs(result)
}

func (s Service) UpdateLatestPosts() {
	result := CallGoroutineChannel()
	savedLatestDate := s.repo.GetLatestPost()
	// savedLatestDate := time.Date(2000, time.August, 7, 0, 0, 0, 0, time.UTC) // sample for force update
	var filterResult []Post
	for _, res := range result {
		if res.Date.After(savedLatestDate) {
			filterResult = append(filterResult, res)
		}
	}
	s.repo.InsertBlogs(filterResult)
}

func (s *Service) GetPostsByCorp(corp string) []dto.PostDTO {
	queryResult := s.repo.SearchBlogsByCrop(corp)
	var result = make([]dto.PostDTO, len(queryResult))
	for idx, post := range queryResult {
		result[idx] = dto.ConvertToDTO(post)
	}
	return result
}

func (s *Service) SearchPosts(value string) []dto.PostDTO {
	queryResult := s.repo.SearchBlogs(value)
	var result = make([]dto.PostDTO, len(queryResult))
	for idx, post := range queryResult {
		result[idx] = dto.ConvertToDTO(post)
	}
	return result
}

func (s *Service) UpdatePostContent(posts []dto.PostDTO) {
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup

	for _, post := range posts {
		wg.Add(1)
		semaphore <- struct{}{} // 세마포어를 획득

		go func(post dto.PostDTO) {
			defer wg.Done()
			defer func() { <-semaphore }() // 작업이 끝나면 세마포어를 해제

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", post.Url, nil)
			if err != nil {
				log.Printf("Error creating request for post %s: %v", post.Id, err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Error fetching URL for post %s: %v \n %v", post.Id, post.Url, err)
				return
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("Error parsing HTML for post %s: %v", post.Id, err)
				return
			}

			doc.Find("script, style").Remove()
			text := strings.TrimSpace(doc.Text())
			s.UpdatePost(post.Id, text)
		}(post)
	}

	wg.Wait() // 모든 고루틴이 완료될 때까지 대기
	log.Println("All posts have been processed")
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
		bucketplace.NewBucketplace().CallApi,
		line.NewLine().CallApi,
		naverpay.NewNaverpay().CallApi,
		socar.NewSocar().CallApi,
		kakaopay.NewKakaopay().CallApi,
		kakaobank.NewKakaobank().CallApi,
		oliveyoung.NewOliveyoung().CallApi,
		daangn.NewDaangn().CallApi,
		toss.NewToss().CallApi,
		musinsa.NewMusinsa().CallApi,
		twonine.NewTwonine().CallApi,
		buzzvil.NewBuzzvil().CallApi,
		kurly.NewKurly().CallApi,
		devsisters.NewDevsisters().CallApi,
		woowa.NewWoowa().CallApi,
	}
	resultChan := make(chan []Post, len(scrapers))

	var wg sync.WaitGroup
	for _, scraper := range scrapers {
		wg.Add(1)
		go func(scrapingFunc func() []Post) {
			defer wg.Done()
			resultChan <- scrapingFunc()
		}(scraper)
	}

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

func (s *Service) UpdateDate() {
	s.repo.UpdateDate()
}
