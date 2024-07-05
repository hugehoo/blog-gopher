package main

import (
	"blog-gopher/scrapper/kakaopay"
	"blog-gopher/scrapper/oliveyoung"
)

type post struct {
	title   string
	url     string
	summary string
	date    string
}

func main() {
	kakaopay.Main()
	oliveyoung.Main()
}
