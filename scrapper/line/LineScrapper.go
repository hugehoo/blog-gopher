package line

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
)

type Line struct {
}

func NewLine() *Line {
	return &Line{}
}

func (l *Line) CallApi() []Post {
	return l.GetPages(1)
}

const baseUrl = "https://techblog.lycorp.co.jp/ko/"

func (l *Line) GetPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = "https://techblog.lycorp.co.jp/page-data/ko/page-data.json"
	} else {
		url = "https://techblog.lycorp.co.jp/page-data/ko/page/" + strconv.Itoa(page) + "/page-data.json"
	}
	res, err := http.Get(url)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	if CheckCodeNonFatal(res) != nil {
		return posts
	}
	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if CheckErrNonFatal(err) != nil {
		return posts
	}
	for _, data := range response.Result.Data.LatestBlog.Edges {
		parsedTime, _ := time.Parse(time.RFC3339Nano, data.Node.PubDate)
		post := Post{
			Title:   data.Node.Title,
			Url:     baseUrl + data.Node.Slug,
			Summary: "",
			Date:    parsedTime,
			Corp:    company.LINE}
		posts = append(posts, post)
	}
	return posts
}

type Response struct {
	Result struct {
		Data struct {
			LatestBlog struct {
				Edges []struct {
					Node struct {
						Slug    string `json:"slug"`
						Title   string `json:"title"`
						PubDate string `json:"pubdate"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"latestBlog"`
		} `json:"data"`
	} `json:"result"`
}
