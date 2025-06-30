package bucketplace

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Bucketplace struct {
}

func (b *Bucketplace) CallApi() []Post {
	resultChan := make(chan []Post)
	var wg sync.WaitGroup

	maxPages := 5
	for i := 1; i < maxPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pages := b.GetPages(page)
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

func (b *Bucketplace) GetPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = "https://www.bucketplace.com/page-data/culture/Tech/page-data.json"
	} else {
		url = "https://www.bucketplace.com/page-data/culture/Tech/" + strconv.Itoa(page) + "/page-data.json"
	}
	res, err := http.Get(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	for _, data := range response.Result.Data.Posts.Nodes {
		post := Post{
			Title:   data.Frontmatter.Title,
			Url:     "https://www.bucketplace.com" + data.Frontmatter.Slug,
			Summary: data.Frontmatter.Description,
			Date:    data.Frontmatter.Date,
			Corp:    company.OHOUSE}
		posts = append(posts, post)
	}
	return posts
}

type Response struct {
	Result Result `json:"result"`
}

type Result struct {
	Data Data `json:"data"`
}

type Data struct {
	Posts Posts `json:"posts"`
}

type Posts struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Frontmatter Frontmatter `json:"frontmatter"`
	Fields      Fields      `json:"fields"`
}

type Frontmatter struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Date           time.Time `json:"date"`
	AuthorName     string    `json:"authorName"`
	ThumbnailImage string    `json:"thumbnailImage"`
}

type Fields struct {
	Tags []string `json:"tags"`
}
