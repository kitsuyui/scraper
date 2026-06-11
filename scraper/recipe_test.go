package scraper

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
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
		t.Errorf("Must be error when invalid css selector is passed.")
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

func TestXPathScalarResults(t *testing.T) {
	r := &recipes{
		{
			Type:  "xpath",
			Query: "count(//p)",
			Label: "paragraph-count",
		}, {
			Type:  "xpath",
			Query: "string(//h1)",
			Label: "heading",
		}, {
			Type:  "xpath",
			Query: "boolean(//p)",
			Label: "has-paragraph",
		},
	}
	cr, err := r.compile()
	if err != nil {
		t.Fatalf("Must not be error on this recipe.")
	}
	input := bytes.NewBufferString(`
<html>
  <body>
    <h1>Heading</h1>
    <p>first</p>
    <p>second</p>
  </body>
</html>
`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Fatalf("This is valid HTML File")
	}
	if len(results) != len(*r) {
		t.Fatalf("Not match size results and recipes")
	}
	expects := []string{"2", "Heading", "true"}
	for i, expect := range expects {
		if results[i].results.PlainResult == nil {
			t.Errorf("PlainResult is nil for %s", (*r)[i].Type)
			continue
		}
		actual := *results[i].results.PlainResult
		if len(actual) != 1 {
			t.Errorf("Not match result count, expect 1 != result %d", len(actual))
			continue
		}
		if actual[0] != expect {
			t.Errorf("Not match scalar XPath result, expect %s != result %s", expect, actual[0])
		}
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
		t.Errorf("Must be error when invalid css selector is passed.")
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

func TestNestedTableNoLeakage(t *testing.T) {
	r := &recipes{{
		Type:  "table-xpath",
		Query: "//table[not(ancestor::table)]",
		Label: "outer",
	}, {
		Type:  "table-css",
		Query: "table",
		Label: "outer",
	}}
	cr, err := r.compile()
	if err != nil {
		t.Errorf("Must not be error on this recipe.")
	}
	input := bytes.NewBufferString(`
<html>
  <body>
    <table>
      <tr>
        <td><table><tr><td>inner</td></tr></table></td>
        <td>outer</td>
      </tr>
    </table>
  </body>
</html>
`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("This is valid HTML File")
	}
	for k, testCase := range results {
		for _, rows := range *testCase.results.TableResult {
			if len(rows) != 1 {
				t.Errorf("Expected 1 row in outer table (no nested-table row leakage), got %d in %s", len(rows), (*r)[k].Type)
			}
		}
	}
}

func TestTableHeaderCells(t *testing.T) {
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
	<body>
		<table border=1>
		  <tr><th>Name</th><th>Score</th></tr>
		  <tr><th>Alice</th><td>10</td></tr>
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
	expects := [][]string{
		{"Name", "Score"},
		{"Alice", "10"},
	}
	for k, testCase := range results {
		if len(*testCase.results.TableResult) != 1 {
			t.Errorf("Not match size table results in %s", (*r)[k].Type)
			continue
		}
		for _, rows := range *testCase.results.TableResult {
			if len(rows) != len(expects) {
				t.Errorf("Not match size table rows in %s", (*r)[k].Type)
				continue
			}
			for i, expectRow := range expects {
				if len(rows[i]) != len(expectRow) {
					t.Errorf("Not match size table columns in %s", (*r)[k].Type)
					continue
				}
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

func TestMarshalJSONNilResults(t *testing.T) {
	sr := &scrapeResult{
		recipe:  recipe{Type: "css", Query: "title", Label: "title"},
		results: nil,
	}
	_, err := sr.MarshalJSON()
	if err == nil {
		t.Errorf("MarshalJSON must return error when results is nil")
	}
}

func TestMarshalJSONBothResultsNil(t *testing.T) {
	sr := &scrapeResult{
		recipe:  recipe{Type: "css", Query: "title", Label: "title"},
		results: &extractResult{PlainResult: nil, TableResult: nil},
	}
	_, err := sr.MarshalJSON()
	if err == nil {
		t.Errorf("MarshalJSON must return error when both PlainResult and TableResult are nil")
	}
}

func TestInvalidColspanRowspanFallsBackToOne(t *testing.T) {
	r := &recipes{{
		Type:  "table-xpath",
		Query: "//table",
		Label: "table1",
	}}
	cr, err := r.compile()
	if err != nil {
		t.Errorf("Must not be error on this recipe.")
	}
	input := bytes.NewBufferString(`<html><body>
		<table>
			<tr><td colspan="0">a</td><td>b</td></tr>
			<tr><td rowspan="-3">c</td><td>d</td></tr>
		</table>
	</body></html>`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("Must not error on invalid colspan/rowspan lower bounds")
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	tables := *results[0].results.TableResult
	if len(tables) != 1 {
		t.Fatalf("Expected 1 table result, got %d", len(tables))
	}
	rows := tables[0]
	expected := [][]string{
		{"a", "b"},
		{"c", "d"},
	}
	if len(rows) != len(expected) {
		t.Fatalf("Expected %d rows, got %d: %#v", len(expected), len(rows), rows)
	}
	for i := range expected {
		if len(rows[i]) != len(expected[i]) {
			t.Fatalf("Expected %d columns in row %d, got %d: %#v", len(expected[i]), i, len(rows[i]), rows[i])
		}
		for j := range expected[i] {
			if rows[i][j] != expected[i][j] {
				t.Errorf("Expected cell[%d][%d] to be %q, got %q", i, j, expected[i][j], rows[i][j])
			}
		}
	}
}

func TestLargeColspanRowspanIsCapped(t *testing.T) {
	r := &recipes{{
		Type:  "table-xpath",
		Query: "//table",
		Label: "table1",
	}}
	cr, err := r.compile()
	if err != nil {
		t.Errorf("Must not be error on this recipe.")
	}
	// colspan=10000 and rowspan=10000 would create 10^8 map entries without the cap;
	// with maxSpan=100 it creates at most 100*100=10000 entries.
	input := bytes.NewBufferString(`<html><body>
		<table><tr><td colspan="10000" rowspan="10000">x</td></tr></table>
	</body></html>`)
	results, err := cr.extractAll(input)
	if err != nil {
		t.Errorf("Must not error on large colspan/rowspan")
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	rows := *results[0].results.TableResult
	if len(rows) == 0 {
		t.Errorf("Expected non-empty table result")
	}
	// The cell value should be present and capped to maxSpan rows/cols
	if len(rows) > maxSpan {
		t.Errorf("rowspan should be capped: got %d rows, expected <= %d", len(rows), maxSpan)
	}
	if len(rows[0]) > maxSpan {
		t.Errorf("colspan should be capped: got %d cols, expected <= %d", len(rows[0]), maxSpan)
	}
}

func TestZeroMatchResultsAsEmptyArrayNotNull(t *testing.T) {
	types := []struct {
		recipeType string
		query      string
	}{
		{TypeCSS, "nonexistent-tag"},
		{TypeXPath, "//nonexistent-tag"},
		{TypeRegex, "DEFINITELY_NOT_IN_HTML"},
	}
	html := bytes.NewBufferString(`<html><body><p>hello</p></body></html>`)
	buf, _ := io.ReadAll(html)

	for _, tt := range types {
		r := &recipes{{Type: tt.recipeType, Query: tt.query, Label: "test"}}
		cr, err := r.compile()
		if err != nil {
			t.Fatalf("compile failed for %s: %v", tt.recipeType, err)
		}
		results, err := cr.extractAll(bytes.NewReader(buf))
		if err != nil {
			t.Fatalf("extractAll failed for %s: %v", tt.recipeType, err)
		}
		jsonBytes, err := json.Marshal(&results[0])
		if err != nil {
			t.Fatalf("MarshalJSON failed for %s: %v", tt.recipeType, err)
		}
		jsonStr := string(jsonBytes)
		if strings.Contains(jsonStr, `"results":null`) {
			t.Errorf("%s: zero-match result marshaled as null, want []; got %s", tt.recipeType, jsonStr)
		}
		if !strings.Contains(jsonStr, `"results":[]`) {
			t.Errorf("%s: zero-match result did not marshal as []; got %s", tt.recipeType, jsonStr)
		}
	}
}
