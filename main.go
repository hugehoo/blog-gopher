package main

import (
	. "blog-gopher/common/types"
	"blog-gopher/scrapper/banksalad"
	"blog-gopher/scrapper/daangn"
	"blog-gopher/scrapper/kakaopay"
	"blog-gopher/scrapper/oliveyoung"
	"blog-gopher/scrapper/toss"
	"log"
)

func main() {

	var result []Post
	result = append(result, kakaopay.Main()...)
	result = append(result, oliveyoung.Main()...)
	result = append(result, daangn.Main()...)
	result = append(result, toss.Main()...)
	result = append(result, banksalad.Main()...)

	total := len(result)
	for i := 0; i < total; i++ {
		post := result[i]
		log.Println(post.Corp, "|", post.Title)
	}
	log.Println("Total :", total)
}
