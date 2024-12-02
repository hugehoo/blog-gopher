package service

import (
	"blog-gopher/common/dto"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	// 기본 만료시간 5분, 만료된 항목 정리 주기 10분으로 설정
	memoryCache = cache.New(24*time.Hour, 24*time.Hour)
)

// 캐시 헬퍼 함수들
type CacheService struct {
	cache *cache.Cache
}

func NewCacheService() *CacheService {
	return &CacheService{
		cache: memoryCache,
	}
}

func (s *CacheService) Get(key string) (interface{}, bool) {
	return s.cache.Get(key)
}

func (s *CacheService) GetPosts(key string) ([]dto.PostDTO, bool) {
	data, found := s.cache.Get(key)
	if !found {
		return nil, false
	}

	posts, ok := data.([]dto.PostDTO)
	if !ok {
		return nil, false
	}

	return posts, true
}

func (s *CacheService) SetPosts(key string, results []dto.PostDTO) {
	s.cache.Set(key, results, 24*time.Hour)
}

func (s *CacheService) GetPostsByCorp(key string) ([]dto.PostDTO, bool) {
	data, found := s.cache.Get(key)
	if !found {
		return nil, false
	}

	posts, ok := data.([]dto.PostDTO)
	if !ok {
		return nil, false
	}
	return posts, true

}

func (s *CacheService) Set(key string, value interface{}, duration time.Duration) {
	s.cache.Set(key, value, duration)
}

func (s *CacheService) Delete(key string) {
	s.cache.Delete(key)
}

// Gin 미들웨어 예제
func CacheMiddleware(cacheService *CacheService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 캐시 키 생성 (예: URL path)
		key := c.Request.URL.Path

		// 캐시된 데이터 확인
		if data, found := cacheService.Get(key); found {
			c.JSON(200, data)
			c.Abort()
			return
		}

		// 원래 요청 처리를 위해 다음 핸들러로 진행
		c.Next()
	}
}
