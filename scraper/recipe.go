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

// Valid values for the recipe Type field.
const (
	TypeXPath      = "xpath"
	TypeCSS        = "css"
	TypeTableXPath = "table-xpath"
	TypeTableCSS   = "table-css"
	TypeRegex      = "regex"
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

// compiledRecipe holds a compiled recipe ready for extraction.
// Invariant: exactly one of domScraper and textScraper is non-nil.
// This invariant is established by compile() and must be preserved by any
// future constructors. extractAll panics if both are nil.
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

// extractResult holds the output of a single scraper execution.
// Invariant: exactly one of PlainResult and TableResult is non-nil.
// TableResult is set by table-* scrapers; PlainResult by all others.
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
	if r.Type == TypeXPath {
		domScraper, err := xPathScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == TypeCSS {
		domScraper, err := cssSelectorScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == TypeTableXPath {
		domScraper, err := tableXPathScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == TypeTableCSS {
		domScraper, err := tableCSSSelectorScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:     r,
			domScraper: domScraper,
		}, nil
	} else if r.Type == TypeRegex {
		textScraper, err := regularExpressionScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:      r,
			textScraper: textScraper,
		}, nil
	}
	return nil, fmt.Errorf("unimplemented recipe type %q; valid types: %s, %s, %s, %s, %s", r.Type, TypeXPath, TypeCSS, TypeTableXPath, TypeTableCSS, TypeRegex)
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
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	var ers []scrapeResult
	for _, cr := range crs {
		if cr.textScraper != nil {
			er := &scrapeResult{recipe: cr.recipe, results: cr.textScraper.extractFromText(htmlText)}
			ers = append(ers, *er)
		} else if cr.domScraper != nil {
			er := &scrapeResult{recipe: cr.recipe, results: cr.domScraper.extractFromNode(doc)}
			ers = append(ers, *er)
		} else {
			// Both scrapers nil: this violates the compiledRecipe invariant.
			// compile() always sets exactly one; reaching here means a future
			// constructor broke the invariant, and silent skipping would hide it.
			panic(fmt.Sprintf("compiledRecipe invariant violation: both domScraper and textScraper are nil for recipe type %q", cr.recipe.Type))
		}
	}
	return ers, nil
}
