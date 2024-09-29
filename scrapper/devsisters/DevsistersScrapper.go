package devsisters

import (
	company "blog-gopher/common/enum"
	. "blog-gopher/common/response"
	. "blog-gopher/common/types"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://tech.devsisters.com"

// 정리 필요 : 기존 goquery 방식으로 안됨. application/json 형태로 리턴받는 페이지라 그런듯.
// 즉 기존 페이징 유알엘로 접근하면, 해당 페이지의 앞단에서 서버로 별도로 요청을 보냄. 그걸 페치해서 json 형태로 가공하여 앞단에 뿌리는듯.
// 첨엔 고루틴에서 상태 공유가 이상하게 되는바람에 모든 값이 똑같이 나오는줄 알앗는데 그건 아녔다.
// 모든 페이지네이션 유알엘이 결국 동일하게 처리되고 있는 이유는 결국 모르겠는데;;
// 기존 페이지네이션 유알엘은 html 을 반환하지 않았기에 내가 처리할 고쿼리 html 태그가 없었다.
func CallApi() []Post {
	return getJsonPage(0)
}

func getJsonPage(i int) []Post {
	requestUrl := baseURL + "/page-data/index/page-data.json?page=" + strconv.Itoa(i)
	res, err := http.Get(requestUrl)
	defer res.Body.Close()

	//var jsonData map[string]interface{}
	var response JSONResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	CheckErr(err)
	//err = json.Unmarshal([]byte(jsonData), &response)
	var posts []Post
	for _, res := range response.Result.Data.AllMarkdownRemark.Edges {
		date := res.Node.Fields.Date
		path := res.Node.Fields.Path
		parsedDate, _ := time.Parse("2006-01-02", strings.TrimSpace(date)) // 문자열을 날짜로 파싱
		post := Post{
			Title:   res.Node.Frontmatter.Title,
			Url:     baseURL + path,
			Summary: res.Node.Frontmatter.Summary,
			Date:    parsedDate.String(),
			Corp:    company.DEVSISTERS}
		posts = append(posts, post)
	}
	return posts
}

type JSONResponse struct {
	ComponentChunkName string `json:"componentChunkName"`
	Path               string `json:"path"`
	Result             Result `json:"result"`
}

type Result struct {
	Data Data `json:"data"`
}

type Data struct {
	AllMarkdownRemark AllMarkdownRemark `json:"allMarkdownRemark"`
}

type AllMarkdownRemark struct {
	Edges []Edge `json:"edges"`
}

type Edge struct {
	Node Node `json:"node"`
}

type Node struct {
	Frontmatter Frontmatter `json:"frontmatter"`
	Fields      Fields      `json:"fields"`
}

type Frontmatter struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type Fields struct {
	Date    string   `json:"date"`
	Path    string   `json:"path"`
	Authors []Author `json:"authors"`
}

type Author struct {
	Author      string `json:"author"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}
