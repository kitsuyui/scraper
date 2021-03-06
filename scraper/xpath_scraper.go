package scraper

import (
	"github.com/antchfx/xpath"
	htmlquery "github.com/antchfx/xquery/html"
	"golang.org/x/net/html"
)

type xPathScraper struct {
	query *xpath.Expr
}

func xPathScraperFromQuery(query string) (*xPathScraper, error) {
	expr, err := xpath.Compile(query)
	if err != nil {
		return nil, err
	}
	return &xPathScraper{query: expr}, nil
}

func (x xPathScraper) extractFromNode(n *html.Node) *extractResult {
	nav := htmlquery.CreateXPathNavigator(n)
	iter := x.query.Evaluate(nav).(*xpath.NodeIterator)
	var ret []string
	for iter.MoveNext() {
		ret = append(ret, iter.Current().Value())
	}
	return &extractResult{PlainResult: &ret}
}
