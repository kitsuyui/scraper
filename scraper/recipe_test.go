package scraper

import (
	"bytes"
	"testing"
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

func TestTableXPathInvalid(t *testing.T) {
	_, err := (&recipe{
		Type:  "table-xpath",
		Query: "!!!!title", // Invalid XPath
		Label: "title",
	}).compile()
	if err == nil {
		t.Errorf("Must be error when invalid xpath is passed.")
	}
}

func TestTableCSSInvalid(t *testing.T) {
	_, err := (&recipe{
		Type:  "table-css",
		Query: "!!!!title", // Invalid CSS Selector
		Label: "title",
	}).compile()
	if err == nil {
		t.Errorf("Must be error when invalid css selecter is passed.")
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
		}, {
			Type:  "regex",
			Query: "Cat, (.*?), Snake",
			Label: "regextest",
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
		<dummy>Cat, Dog, Snake</dummy>
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
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("This is valid HTML File")
	}
	if len(results) != len(*r) {
		t.Errorf("Not match size results and recipes")
	}
	if (*(results[0].results.PlainResult))[0] != "test passed" {
		t.Errorf("Invalid")
	}
	if (*(results[1].results.PlainResult))[0] != "test passed" {
		t.Errorf("Invalid")
	}
	if (*(results[4].results.PlainResult))[0] != "Cat, Dog, Snake" {
		t.Errorf("Invalid")
	}
}

func TestBasicsTables(t *testing.T) {
	r := &recipes{{
		Type:  "table-xpath",
		Query: "//table",
		Label: "table1",
	}, {
		Type:  "table-css",
		Query: "table",
		Label: "table1",
	}}
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
		  <tr><td>a</td><td>b</td><td>c</td><td>d</td></tr>
		  <tr><td>e</td><td rowspan="2" colspan="2">f</td><td>g</td></tr>
		  <tr><td>h</td><td>i</td></tr>
		</table>
	</body>
</html>
`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("This is valid HTML File")
	}
	if len(results) != len(*r) {
		t.Errorf("Not match size results and recipes")
	}
	for k, testCase := range results {
		for _, rows := range *testCase.results.TableResult {
			if len(rows) != 3 {
				t.Errorf("Not match size table rows")
			}
			for _, row := range rows {
				if len(row) != 4 {
					t.Errorf("Not match size table columns")
				}
			}
			expects := [][]string{
				{"a", "b", "c", "d"},
				{"e", "f", "f", "g"},
				{"h", "f", "f", "i"},
			}
			for i, expectRow := range expects {
				for j, cell := range expectRow {
					if rows[i][j] != cell {
						t.Errorf("Not match table cell, expect %s != result %s in %s", cell, rows[i][j], (*r)[k].Type)
					}
				}
			}
		}
	}
}

func TestInvalidRowspanColspan(t *testing.T) {
	r := &recipes{
		{
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
	<table border=1>
	  <tr><td>a</td><td>b</td><td>c</td><td>d</td></tr>
	  <tr><td>e</td><td rowspan="a" colspan="b">f</td><td>g</td></tr>
	  <tr><td>h</td><td>i</td></tr>
	</table>
</html>
`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("This is valid HTML File")
	}
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

func TestRegexInvalid(t *testing.T) {
	_, err := (&recipe{
		Type:  "regex",
		Query: "(.*?", // Invalid CSS Selector
		Label: "something",
	}).compile()
	if err == nil {
		t.Errorf("Must be error when invalid regular expression is passed.")
	}
}
