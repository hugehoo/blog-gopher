package kmong

import (
	. "blog-gopher/common/types"
	"blog-gopher/common/utils"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

var urls = []string{
	"https://blog.kmong.com/tech/home",
}

func CallApi() []Post {
	return utils.CallGoroutineApi(getPages, urls)
}

func getPages(page string) []Post {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// Add timeout to prevent infinite execution
	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var posts []Post

	// create channels for communication
	postChan := make(chan Post)
	done := make(chan bool)

	// collect posts in a separate goroutine
	go func() {
		for post := range postChan {
			posts = append(posts, post)
		}
		done <- true
	}()

	// navigate and scrape
	err := chromedp.Run(ctx,
		chromedp.Navigate(page),
		chromedp.WaitVisible(".streamItem.streamItem--section"),
		scrapePostsWithScroll(postChan),
	)
	if err != nil {
		log.Fatal(err)
	}

	close(postChan)
	<-done

	// print results
	fmt.Printf("Scraped %d blog posts\n", len(posts))
	for i, p := range posts {
		fmt.Printf("Post %d:\n", i+1)
		fmt.Printf("Title: %s\nSubtitle",
			p.Title)
	}

	return nil
}

func scrapePostsWithScroll(postChan chan<- Post) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var previousHeight int
		var processedURLs = make(map[string]bool)

		for {
			// get all post nodes within the main container
			var nodes []*cdp.Node
			if err := chromedp.Nodes(".col.u-xs-size12of12.js-trackPostPresentation", &nodes).Do(ctx); err != nil {
				return err
			}

			// extract data from each post
			for _, node := range nodes {
				var post Post
				var postURL string

				// Get the post URL first to check if we've already processed this post
				if err := chromedp.AttributeValue("a[data-action=\"open-post\"]", "href", &postURL, nil, chromedp.ByQuery, chromedp.FromNode(node)); err != nil {
					continue
				}

				// Skip if we've already processed this post
				if processedURLs[postURL] {
					continue
				}

				// Mark this URL as processed
				processedURLs[postURL] = true
				post.Url = postURL

				// Extract other post data
				if err := chromedp.Run(ctx,
					chromedp.Text(".u-fontSize24", &post.Title, chromedp.ByQuery, chromedp.FromNode(node)),
					//chromedp.Text(".u-fontSize18", &post.Subtitle, chromedp.ByQuery, chromedp.FromNode(node)),
					//chromedp.Text(".ds-link--styleSubtle", &post.Author, chromedp.ByQuery, chromedp.FromNode(node)),
					chromedp.Text("time", &post.Date, chromedp.ByQuery, chromedp.FromNode(node)),
				); err != nil {
					continue
				}

				// Extract image URL from background-image style
				var backgroundImage string
				if err := chromedp.AttributeValue(".u-block", "style", &backgroundImage, nil, chromedp.ByQuery, chromedp.FromNode(node)); err == nil {
					//if len(backgroundImage) > 0 {
					//	post.ImageURL = extractImageURL(backgroundImage)
					//}
				}

				// send post to channel if title is not empty
				if post.Title != "" {
					postChan <- post
				}
			}

			// Check scroll height
			var height int
			if err := chromedp.Evaluate(`document.documentElement.scrollHeight`, &height).Do(ctx); err != nil {
				return err
			}

			// Break if we've reached the bottom (no height change after scroll)
			if height == previousHeight {
				// Try one more time to ensure we're really at the bottom
				if err := chromedp.Run(ctx,
					chromedp.Sleep(2*time.Second),
				); err != nil {
					return err
				}

				// Check height again
				var newHeight int
				if err := chromedp.Evaluate(`document.documentElement.scrollHeight`, &newHeight).Do(ctx); err != nil {
					return err
				}

				if newHeight == height {
					break
				}
			}
			previousHeight = height

			// Scroll and wait for content to load
			if err := chromedp.Run(ctx,
				chromedp.Evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil),
				chromedp.Sleep(2*time.Second), // Wait for new content to load
			); err != nil {
				return err
			}

			// Optional: Add a progress indicator
			log.Printf("Scrolled to height: %d, Found %d posts so far\n", height, len(processedURLs))
		}

		return nil
	}
}

func extractImageURL(style string) string {
	start := strings.Index(style, "url(\"") + 5
	end := strings.Index(style, "\")")
	if start >= 5 && end > start {
		return style[start:end]
	}
	return ""
}
