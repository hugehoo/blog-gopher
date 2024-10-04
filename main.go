package main

import (
	"blog-gopher/common/dto"
	"blog-gopher/common/types"
	"blog-gopher/repository"
	"blog-gopher/service"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Println("✅ Starting initialization...")
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//
	var s *service.Service
	if repository, err := repository.NewRepository(); err != nil {
		panic(err)
	} else {
		s = service.NewService(repository)
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	r.GET("/posts", func(c *gin.Context) {
		posts := getPosts(c, s)
		c.JSON(http.StatusOK, posts)
	})

	r.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, searchPosts(c, s))
	})

	ginLambda = ginadapter.New(r)
	log.Println("✅ Initialization complete")
}

func main() {
	log.Println("📌 Start Lambda")
	lambda.Start(Handler)
	//localHandler()
}

func localHandler() {
	log.Println("📌 Start Local")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//
	var s *service.Service
	if repository, err := repository.NewRepository(); err != nil {
		panic(err)
	} else {
		s = service.NewService(repository)
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	r.GET("/posts", func(c *gin.Context) {
		posts := getPosts(c, s)
		c.JSON(http.StatusOK, posts)
	})

	r.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, searchPosts(c, s))
	})

	r.POST("/", func(c *gin.Context) {
		s.UpdateAllPosts()
	})

	r.POST("/content", func(c *gin.Context) {
		updatePostContent(c, s)
	})

	r.Run()
}

func updatePostContent(c *gin.Context, s *service.Service) {
	var corp []string
	posts := s.GetPosts(corp, 0, 1)
	maxConcurrent := 10
	semaphore := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup

	for _, post := range posts {
		wg.Add(1)
		semaphore <- struct{}{} // 세마포어를 획득

		go func(post types.Post) {
			defer wg.Done()
			defer func() { <-semaphore }() // 작업이 끝나면 세마포어를 해제

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", post.Url, nil)
			if err != nil {
				log.Printf("Error creating request for post %s: %v", post.ID, err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Error fetching URL for post %s: %v \n %v", post.ID, post.Url, err)
				return
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("Error parsing HTML for post %s: %v", post.ID, err)
				return
			}

			doc.Find("script, style").Remove()
			text := strings.TrimSpace(doc.Text())
			s.UpdatePost(post.ID, text)
		}(post)
	}

	wg.Wait() // 모든 고루틴이 완료될 때까지 대기
	log.Println("All posts have been processed")
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling request: %s %s", request.HTTPMethod, request.Path)
	return ginLambda.ProxyWithContext(ctx, request)
}

func searchPosts(c *gin.Context, s *service.Service) []dto.PostDTO {
	pageStr := c.Query("page")         // 쿼리 파라미터에서 페이지 번호를 가져옵니다.
	pageSizeStr := c.Query("pageSize") // 쿼리 파라미터에서 페이지 크기를 가져옵니다.

	// 기본 값 설정
	page := 1
	pageSize := 20

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
	value := c.Query("keyword")
	var queryResult []types.Post
	queryResult = s.SearchPosts(value, page, pageSize)

	var result = make([]dto.PostDTO, len(queryResult))
	for idx, post := range queryResult {
		result[idx] = dto.ConvertToDTO(post)
	}
	return result
}

type Category struct {
	Category string `form:"category"`
}

func getPosts(c *gin.Context, s *service.Service) []dto.PostDTO {
	pageStr := c.Query("page")         // 쿼리 파라미터에서 페이지 번호를 가져옵니다.
	pageSizeStr := c.Query("pageSize") // 쿼리 파라미터에서 페이지 크기를 가져옵니다.

	// 기본 값 설정
	page := 1
	pageSize := 20

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

	results := s.GetPosts(splits, page, pageSize)
	response := make([]dto.PostDTO, len(results))
	for idx, result := range results {
		response[idx] = dto.ConvertToDTO(result)
	}
	return response
}
