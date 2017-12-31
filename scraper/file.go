package scraper

import (
	"encoding/json"
	"io"
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
	ers, err := cr.extractAll(input)
	if err != nil {
		return err
	}
	e := json.NewEncoder(output)
	e.SetIndent(" ", "  ")
	e.Encode(ers)
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
