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
			Query: "//title",
			Label: "title",
		}, {
			Type:  "css",
			Query: "title",
			Label: "title",
		}, {
			Type:  "table-css",
			Query: "table",
			Label: "table1",
		}, {
			Type:  "table-xpath",
			Query: "//table",
			Label: "table1",
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
	<body>
		<table border=1>
		  <tr><td>a</td><td>b</td><td>c</td></tr>
		  <tr><td>d</td><td>e</td><td>f</td></tr>
		</table>
		<hr>
		<table border=1>
		  <tr><td>a</td><td>b</td><td>c</td><td rowspan="2">d</td></tr>
		  <tr><td>e</td><td colspan="2">f</td></tr>
		  <tr><td>i</td><td>j</td><td>k</td><td>l</td></tr>
		</table>
		<hr>
		<table border=1>
		  <tr><td>a</td><td>b</td><td rowspan="2" bgcolor="pink">c</td><td>d</td></tr>
		  <tr><td>e</td><td colspan="3" bgcolor="yellow">f</td></tr>
		  <tr><td>i</td><td>j</td><td>k</td><td>l</td></tr>
		</table>
		<hr>
		<table border=1>
		  <tr><td>a</td><td>b</td><td>c</td><td>d</td></tr>
		  <tr><td>e</td><td rowspan="2" colspan="2">f</td><td>g</td></tr>
		  <tr><td>h</td><td>i</td></tr>
		</table>
	</body>
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
