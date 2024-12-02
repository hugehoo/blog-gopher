package types

import (
	compnay "blog-gopher/common/enum"
	"time"
)

type Post struct {
	ID      string `bson:"_id,omitempty"`
	Title   string
	Url     string
	Summary string
	Date    time.Time
	Content string
	Corp    compnay.Company
}
