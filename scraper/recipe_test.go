package scraper

import (
	"bytes"
	"testing"

	htmlquery "github.com/antchfx/xquery/html"
)

func TestNotImplementedType(t *testing.T) {
	_, err := (&recipe{
		Type:  "yay_this_is_unimplemented",
		Query: "***",
		Label: "**",
	}).compile()
	if err == nil {
		t.Errorf("Overlook Unimplemented type.")
	}
}

func TestCSSValid(t *testing.T) {
	_, err := (&recipe{
		Type:  "css",
		Query: "title",
		Label: "title",
	}).compile()
	if err != nil {
		t.Errorf("This is valid CSS Selector.")
	}
}

func TestCSSInvalid(t *testing.T) {
	_, err := (&recipe{
		Type:  "css",
		Query: "<title", // Invalid CSS Selector
		Label: "title",
	}).compile()
	if err == nil {
		t.Errorf("Must be error when invalid css selecter is passed.")
	}
}

func TestXPathValid(t *testing.T) {
	_, err := (&recipe{
		Type:  "xpath",
		Query: "//title/text()",
		Label: "title",
	}).compile()
	if err != nil {
		t.Errorf("This is valid XPath expression.")
	}
}

func TestXPathInvalid(t *testing.T) {
	_, err := (&recipe{
		Type:  "xpath",
		Query: "!!!!title", // Invalid XPath
		Label: "title",
	}).compile()
	if err == nil {
		t.Errorf("Must be error when invalid xpath is passed.")
	}
}

func TestBasics(t *testing.T) {
	r := &recipes{
		{
			Type:  "xpath",
			Query: "//title", // Invalid XPath
			Label: "title",
		}, {
			Type:  "css",
			Query: "title", // Invalid CSS Selector
			Label: "title",
		},
	}
	cr, err := r.compile()
	if err != nil {
		t.Errorf("Must not be error on this recipe.")
	}
	input := bytes.NewBufferString(`
<html>
  <head>
    <title>test passed</title>
  </head>
</html>
`)
	n, err := htmlquery.Parse(input)
	if err != nil {
		t.Errorf("This is valid HTML File")
	}
	results := cr.extractAll(n)
	if len(results) != len(*r) {
		t.Errorf("Not match size results and recipes")
	}
}

func TestInvalidRecipe(t *testing.T) {
	r := &recipes{
		{
			Type:  "xpath",
			Query: "<title", // Invalid XPath
			Label: "title",
		}, {
			Type:  "css",
			Query: "<title", // Invalid XPath
			Label: "title",
		},
	}
	_, err := r.compile()
	if err == nil {
		t.Errorf("Must be error on this recipe.")
	}
}
