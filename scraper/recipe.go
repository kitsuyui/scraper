package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	htmlquery "github.com/antchfx/xquery/html"
	"golang.org/x/net/html"
)

type recipe struct {
	Type  string `json:"type"`
	Label string `json:"label"` // For human readability
	Query string `json:"query"`
}

type scrapeResult struct {
	recipe
	results *extractResult
}

func (sr *scrapeResult) MarshalJSON() ([]byte, error) {
	if sr.results.TableResult != nil {
		return json.Marshal(&struct {
			recipe
			Results [][][]string `json:"results"`
		}{
			recipe:  sr.recipe,
			Results: *sr.results.TableResult,
		})
	}
	return json.Marshal(&struct {
		recipe
		Results []string `json:"results"`
	}{
		recipe:  sr.recipe,
		Results: *sr.results.PlainResult,
	})
}

type recipes []recipe

type compiledRecipe struct {
	domScraper  domScraper
	textScraper textScraper
	recipe
}

type compiledRecipes []compiledRecipe

type textScraper interface {
	extractFromText(string) *extractResult
}

type domScraper interface {
	extractFromNode(*html.Node) *extractResult
}

type extractResult struct {
	PlainResult *[]string
	TableResult *[][][]string
}

func (rs recipes) compile() (compiledRecipes, error) {
	var crs []compiledRecipe
	var errMsgs []string
	for i, r := range rs {
		cr, err := r.compile()
		if err != nil {
			errMsg := fmt.Sprintf("Error: Recipe[%d]: %s", i, err.Error())
			errMsgs = append(errMsgs, errMsg)
		} else {
			crs = append(crs, *cr)
		}
	}
	if len(errMsgs) == 0 {
		return crs, nil
	}
	return crs, fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
}

func (r recipe) compile() (*compiledRecipe, error) {
	if r.Type == "xpath" {
		domScraper, err := xPathScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == "css" {
		domScraper, err := cssSelectorScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == "table-xpath" {
		domScraper, err := tableXPathScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == "table-css" {
		domScraper, err := tableCSSSelectorScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == "regex" {
		textScraper, err := regularExpressionScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:      r,
			textScraper: textScraper,
		}, nil
	}
	return nil, fmt.Errorf("Unimplemented type: %s", r.Type)
}

func (crs compiledRecipes) extractAll(input io.Reader) ([]scrapeResult, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(input)
	if err != nil {
		return nil, err
	}
	htmlText := buf.String()
	input = strings.NewReader(htmlText)
	doc, err := htmlquery.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("If this error is occurred, please tell me HTML for adding unit-test case.%s", err)
	}
	var ers []scrapeResult
	for _, cr := range crs {
		if cr.textScraper != nil {
			er := &scrapeResult{recipe: cr.recipe, results: cr.textScraper.extractFromText(htmlText)}
			ers = append(ers, *er)
		} else if cr.domScraper != nil {
			er := &scrapeResult{recipe: cr.recipe, results: cr.domScraper.extractFromNode(doc)}
			ers = append(ers, *er)
		}
	}
	return ers, nil
}
