package scraper

import (
	"regexp"
)

type regularExpressionScraper struct {
	query *regexp.Regexp
}

func regularExpressionScraperFromQuery(query string) (*regularExpressionScraper, error) {
	expr, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}
	return &regularExpressionScraper{query: expr}, nil
}

func (c regularExpressionScraper) extractFromText(htmlText string) *extractResult {
	ret := c.query.FindStringSubmatch(htmlText)
	return &extractResult{PlainResult: &ret}
}
