package dto

import (
	compnay "blog-gopher/common/enum"
	. "blog-gopher/common/types"
)

type PostDTO struct {
	Id      string          `json:"id"`
	Title   string          `json:"title"`
	Url     string          `json:"url"`
	Summary string          `json:"summary"`
	Date    string          `json:"date"`
	Corp    compnay.Company `json:"corp"`
}

func ConvertToDTO(post Post) PostDTO {
	dateOnly := post.Date.Format("2006-01-02")
	return PostDTO{
		Id:      post.ID,
		Title:   post.Title,
		Url:     post.Url,
		Summary: post.Summary,
		Date:    dateOnly,
		Corp:    post.Corp,
	}
}
