package main

import (
	"blog-gopher/common/dto"
	"blog-gopher/config"
	"blog-gopher/repository"
	"blog-gopher/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var pathFlag = flag.String("config", "./config.toml", "config set up")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)
	var s *service.Service
	if repository, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		s = service.NewService(repository)
	}

	r := gin.Default()
	r.GET("/posts", func(c *gin.Context) {
		posts := getPosts(c, s)
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	})

	r.POST("/", func(c *gin.Context) {
		s.UpdateAllPosts()
	})
	r.Run()
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
