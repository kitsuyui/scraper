package scraper

import (
	"bytes"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type cssSelectorScraper struct {
	query cascadia.Selector
}

func cssSelectorScraperFromQuery(query string) (*cssSelectorScraper, error) {
	expr, err := cascadia.Compile(query)
	if err != nil {
		return nil, err
	}
	return &cssSelectorScraper{query: expr}, nil
}

func (c cssSelectorScraper) extract(n *html.Node) *extractResult {
	var ret []string
	for _, m := range c.query.MatchAll(n) {
		ret = append(ret, extractTextFromNodeRecursively(m))
	}
	return &extractResult{PlainResult: &ret}
}

func extractTextFromNodeRecursively(n *html.Node) string {
	var b bytes.Buffer
	if n.Type == html.TextNode {
		b.WriteString(n.Data)
	}
	if n.FirstChild != nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			b.WriteString(extractTextFromNodeRecursively(c))
		}
	}
	return b.String()
}
