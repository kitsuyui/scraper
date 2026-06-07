package scraper

import (
	"fmt"
	"strconv"

	htmlquery "github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
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
	var ret []string
	switch result := x.query.Evaluate(nav).(type) {
	case *xpath.NodeIterator:
		for result.MoveNext() {
			ret = append(ret, result.Current().Value())
		}
	case string:
		ret = append(ret, result)
	case float64:
		ret = append(ret, strconv.FormatFloat(result, 'f', -1, 64))
	case bool:
		ret = append(ret, strconv.FormatBool(result))
	case nil:
	default:
		ret = append(ret, fmt.Sprint(result))
	}
	return &extractResult{PlainResult: &ret}
}
