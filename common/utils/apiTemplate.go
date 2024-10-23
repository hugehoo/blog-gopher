package utils

import (
	"blog-gopher/common/types"
	"sync"
)

func CallApiTemplate(
	getPages func(url string) []types.Post,
	urls []string,
) []types.Post {

	type Result struct {
		Posts []types.Post
	}
	var response []types.Post
	var wg sync.WaitGroup
	resultChan := make(chan Result, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			pages := getPages(url)
			resultChan <- Result{
				Posts: pages,
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		response = append(response, result.Posts...)
	}
	return response
}
