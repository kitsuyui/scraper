package scraper

import (
	"encoding/json"
	"fmt"
	"io"

	htmlquery "github.com/antchfx/xquery/html"
)

// ValidateConfigFile provides syntax-check for scraper config file.
func ValidateConfigFile(f io.Reader) error {
	_, err := compileConfigFile(f)
	return err
}

// ScrapeByConfFile is the main routine of cli.
func ScrapeByConfFile(confFile io.Reader, input io.Reader, output io.Writer) error {
	cr, err := compileConfigFile(confFile)
	if err != nil {
		return err
	}
	return scrapeByCompiledRecipes(cr, input, output)
}

func scrapeByCompiledRecipes(cr compiledRecipes, input io.Reader, output io.Writer) error {
	doc, err := htmlquery.Parse(input)
	if err != nil {
		return fmt.Errorf("If this error is occurred, please tell me HTML for adding unit-test case.%s", err)
	}
	ers := cr.extractAll(doc)
	json.NewEncoder(output).Encode(ers)
	return nil
}

func compileConfigFile(f io.Reader) (compiledRecipes, error) {
	var r recipes
	err := json.NewDecoder(f).Decode(&r)
	if err != nil {
		return nil, err
	}
	cr, err := r.compile()
	if err != nil {
		return nil, err
	}
	return cr, nil
}
