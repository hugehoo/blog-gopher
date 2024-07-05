package types

import (
	compnay "blog-gopher/common/enum"
)

type Post struct {
	Title   string
	Url     string
	Summary string
	Date    string
	Corp    compnay.Company
}
