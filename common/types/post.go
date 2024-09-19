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
	Corp    compnay.Company
}
