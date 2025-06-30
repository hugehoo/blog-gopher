package scrapper

import "blog-gopher/common/types"

type Page interface {
	GetPages(page int) []types.Post
}
