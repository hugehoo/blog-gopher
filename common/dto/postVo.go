package dto

import (
	"blog-gopher/common/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
)

func InsertPost(post types.Post) primitive.D {
	return bson.D{
		{"title", html.UnescapeString(post.Title)},
		{"url", post.Url},
		{"summary", post.Summary},
		{"date", post.Date},
		{"content", post.Content},
		{"corp", post.Corp},
	}
}
