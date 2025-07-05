package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"blog-gopher/common/dto"
	"blog-gopher/repository"
	"blog-gopher/service"
)

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(Handler)
	//localHandler()
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

// lambda handler
func init() {
	log.Println("✅ Starting initialization...")
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	getPing(r)

	var s *service.Service
	if repo, err := repository.NewRepository(); err != nil {
		panic(err)
	} else {
		s = service.NewService(repo)
	}

	cacheService := service.NewCacheService()

	webConfigurations(r)

	getAllPosts(r, cacheService, s)
	getPostsByCorp(r, cacheService, s)
	searchPosts(r, cacheService, s)

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
	getPing(r)

	mongo, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	s := service.NewService(mongo)
	cacheService := service.NewCacheService()

	webConfigurations(r)

	getAllPosts(r, cacheService, s)
	getPostsByCorp(r, cacheService, s)
	searchPosts(r, cacheService, s)

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
		deleteDefaultPostCache(cacheService)
		s.UpdateAllPosts()
	})

	r.POST("/latest", func(c *gin.Context) {
		deleteDefaultPostCache(cacheService)
		s.UpdateLatestPosts()
	})

	r.POST("/content", func(c *gin.Context) {
		posts := s.GetPosts(c)
		s.UpdatePostContent(posts)
	})

	r.POST("/update-date", func(c *gin.Context) {
		s.UpdateDate()
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}

func deleteDefaultPostCache(cacheService *service.CacheService) {
	key := "posts"
	cacheService.Delete(key)
}

func webConfigurations(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))
}

func getPing(r *gin.Engine) gin.IRoutes {
	return r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func searchPosts(r *gin.Engine, cacheService *service.CacheService, s *service.Service) gin.IRoutes {
	return r.GET("/search", func(c *gin.Context) {
		value := c.Query("keyword")
		if value == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
			return
		}

		// Create cache key for search
		cacheKey := fmt.Sprintf("search:%s", value)

		// Check cache first
		if cached, found := cacheService.GetPosts(cacheKey); found {
			log.Printf("Search cache hit for keyword: %s", value)
			c.JSON(http.StatusOK, cached)
			return
		}

		// If not in cache, perform search
		result, err := s.SearchPosts(value)
		if err != nil {
			log.Printf("Search error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}

		// Cache the results
		cacheService.SetPosts(cacheKey, result)
		c.JSON(http.StatusOK, result)
	})
}

func getPostsByCorp(r *gin.Engine, cacheService *service.CacheService, s *service.Service) gin.IRoutes {
	return r.GET("/corps", func(c *gin.Context) {
		value := c.Query("corp")
		key := fmt.Sprintf("corps::%s", value)
		cache, exists := cacheService.GetPostsByCorp(key)
		if exists {
			log.Println("Hit the cache")
			c.JSON(http.StatusOK, cache)
		} else {
			posts := s.GetPostsByCorp(value)
			cacheService.SetPosts(key, posts)
			c.JSON(http.StatusOK, posts)
		}
	})
}

func getAllPosts(r *gin.Engine, cacheService *service.CacheService, s *service.Service) gin.IRoutes {
	return r.GET("/posts", func(c *gin.Context) {
		key := "posts"
		cache, bool := cacheService.GetPosts(key)
		if bool != false {
			c.JSON(http.StatusOK, cache)
		} else {
			results := s.GetPosts(c)
			cacheService.SetPosts(key, results)
			c.JSON(http.StatusOK, results)
		}
	})
}
