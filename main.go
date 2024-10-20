package main

import (
	"blog-gopher/common/dto"
	"blog-gopher/repository"
	"blog-gopher/service"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var ginLambda *ginadapter.GinLambda

func main() {
	log.Println("📌 Start Lambda")
	lambda.Start(Handler)
	//localHandler()
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Handling request: %s %s", request.HTTPMethod, request.Path)
	return ginLambda.ProxyWithContext(ctx, request)
}

// lambda handler
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
		results := s.GetPosts(c)
		c.JSON(http.StatusOK, results)
	})

	r.GET("/search", func(c *gin.Context) {
		value := c.Query("keyword")
		result := s.SearchPosts(value)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/search/blog/:blogId", func(c *gin.Context) {
		result := s.SearchPostsById(c.Param("blogId"))
		c.JSON(http.StatusOK, dto.ConvertToDTO(result))
	})

	ginLambda = ginadapter.New(r)
	log.Println("✅ Initialization complete")
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
		posts := s.GetPosts(c)
		c.JSON(http.StatusOK, posts)
	})

	r.GET("/search", func(c *gin.Context) {
		value := c.Query("keyword")
		result := s.SearchPosts(value)
		c.JSON(http.StatusOK, result)
	})

	r.GET("/auto-search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		response := s.AutoSearchPosts(keyword)
		c.JSON(http.StatusOK, response)
	})

	r.GET("/search/blog/:blogId", func(c *gin.Context) {
		result := s.SearchPostsById(c.Param("blogId"))
		c.JSON(http.StatusOK, dto.ConvertToDTO(result))
	})

	r.POST("/", func(c *gin.Context) {
		s.UpdateAllPosts()
	})

	r.POST("/latest", func(c *gin.Context) {
		s.UpdateLatestPosts()
	})

	r.POST("/content", func(c *gin.Context) {
		posts := s.GetPosts(c)
		s.UpdatePostContent(posts)
	})

	r.Run()
}
