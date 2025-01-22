package crawl

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func CrawlURL(urlChan chan<- string, errChan chan<- error, link string) {
	resp, err := http.Get(link)
	if err != nil {
		errChan <- fmt.Errorf("failed to perform get-request, err: %w", err)
		return
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("failed to parse html, err: %w", err)
		return
	}

	seen := make(map[string]bool)
	var recursion func(*html.Node)
	recursion = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(attr.Val)
				if err != nil {
					continue
				}
				if link.Hostname() != resp.Request.URL.Hostname() || seen[link.String()] {
					continue
				}
				seen[link.String()] = true
				fmt.Println(link.String())
				urlChan <- link.String()
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			recursion(c)
		}
	}
	
	recursion(root)

	defer close(urlChan)
}
