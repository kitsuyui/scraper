package scraper

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type recipe struct {
	Type  string `json:"type"`
	Label string `json:"label"` // For human readability
	Query string `json:"query"`
}

type extractedResult struct {
	recipe
	Results []string `json:"results"`
}

type recipes []recipe

type compiledRecipe struct {
	scraper scraper
	recipe
}

type compiledRecipes []compiledRecipe

type scraper interface {
	extract(n *html.Node) (ret []string)
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
	return crs, fmt.Errorf(strings.Join(errMsgs, "\n"))
}

func (r recipe) compile() (*compiledRecipe, error) {
	if r.Type == "xpath" {
		scraper, err := xPathScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:  r,
			scraper: scraper,
		}, nil
	} else if r.Type == "css" {
		scraper, err := cssSelectorScraperFromQuery(r.Query)
		if err != nil {
			return nil, err
		}
		return &compiledRecipe{
			recipe:  r,
			scraper: scraper,
		}, nil
	}
	return nil, fmt.Errorf("Unimplemented type: %s", r.Type)
}

// ExtractAll do every recipe.
func (crs compiledRecipes) extractAll(n *html.Node) (ers []extractedResult) {
	for _, cr := range crs {
		e := cr.scraper.extract(n)
		er := &extractedResult{recipe: cr.recipe, Results: e}
		ers = append(ers, *er)
	}
	return ers
}
