package dto

import (
	compnay "blog-gopher/common/enum"
	. "blog-gopher/common/types"
)

type PostDTO struct {
	Title   string          `json:"title"`
	Url     string          `json:"url"`
	Summary string          `json:"summary"`
	Date    string          `json:"date"`
	Corp    compnay.Company `json:"corp"`
}

func ConvertToDTO(post Post) PostDTO {
	return PostDTO{
		Title:   post.Title,
		Url:     post.Url,
		Summary: post.Summary,
		Date:    post.Date,
		Corp:    post.Corp,
	}
}
