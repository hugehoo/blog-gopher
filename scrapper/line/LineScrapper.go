package line

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func CallApi() []Post {
	var result []Post
	for i := 1; i < 6; i++ {
		pages := getPages(i)
		result = append(result, pages...)
	}
	return result
}

const baseUrl = "https://techblog.lycorp.co.jp/ko/"

func getPages(page int) []Post {
	var posts []Post
	var url string
	if page == 1 {
		url = "https://techblog.lycorp.co.jp/page-data/ko/page-data.json"
	} else {
		url = "https://techblog.lycorp.co.jp/page-data/ko/page/" + strconv.Itoa(page) + "/page-data.json"
	}
	res, err := http.Get(url)
	CheckErr(err)
	CheckCode(res)
	defer res.Body.Close()

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	for _, data := range response.Result.Data.BlogsQuery.Edges {
		parsedTime, _ := time.Parse(time.RFC3339Nano, data.Node.PubDate)
		post := Post{
			Title:   data.Node.Title,
			Url:     baseUrl + data.Node.Slug,
			Summary: "",
			Date:    parsedTime.String(),
			Corp:    company.LINE}
		posts = append(posts, post)
	}
	return posts
}

type Response struct {
	Result struct {
		Data struct {
			BlogsQuery struct {
				Edges []struct {
					Node struct {
						Slug    string `json:"slug"`
						Title   string `json:"title"`
						PubDate string `json:"pubdate"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"BlogsQuery"`
		} `json:"data"`
	} `json:"result"`
}
