package types

import (
	compnay "blog-gopher/common/enum"
)

type Post struct {
	ID      string `bson:"_id,omitempty"`
	Title   string
	Url     string
	Summary string
	Date    string
	Content string
	Corp    compnay.Company
}
