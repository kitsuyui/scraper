package scraper

import (
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	htmlquery "github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

type tableXPathScraper struct {
	query *xpath.Expr
}

type tableCSSSelectorScraper struct {
	query cascadia.Selector
}

func tableXPathScraperFromQuery(query string) (*tableXPathScraper, error) {
	expr, err := xpath.Compile(query)
	if err != nil {
		return nil, err
	}
	return &tableXPathScraper{query: expr}, nil
}

func (x tableXPathScraper) extractFromNode(n *html.Node) *extractResult {
	var ret [][][]string
	for _, elem := range htmlquery.Find(n, x.query.String()) {
		k := extractTable(elem)
		ret = append(ret, k)
	}
	return &extractResult{TableResult: &ret}
}

func tableCSSSelectorScraperFromQuery(query string) (*tableCSSSelectorScraper, error) {
	expr, err := cascadia.Compile(query)
	if err != nil {
		return nil, err
	}
	return &tableCSSSelectorScraper{query: expr}, nil
}

func (c tableCSSSelectorScraper) extractFromNode(n *html.Node) *extractResult {
	var ret [][][]string
	for _, m := range c.query.MatchAll(n) {
		ret = append(ret, extractTable(m))
	}
	return &extractResult{TableResult: &ret}
}

func extractTable(n *html.Node) [][]string {
	colmax := 0
	rowmax := 0
	c := map[int]map[int]*string{}
	for i, tr := range htmlquery.Find(n, ".//tr") {
		jFixed := 0
		for _, td := range htmlquery.Find(tr, ".//th or .//td") {
			colspan, rowspan := parseColspanRowspan(td)
			strVal := extractTextFromNodeRecursively(td)
			for isFilled(c, i, jFixed) {
				jFixed++
			}
			editCellInRange(c, i, jFixed, colspan, rowspan, strVal)
			if i+rowspan > rowmax {
				rowmax = i + rowspan
			}
			if jFixed+colspan > colmax {
				colmax = jFixed + colspan
			}
			jFixed += colspan
		}
	}
	return mapTableToSliceTable(c, rowmax, colmax)
}

func mapTableToSliceTable(c map[int]map[int]*string, rowmax int, colmax int) (v [][]string) {
	for i := 0; i < rowmax; i++ {
		var row []string
		for j := 0; j < colmax; j++ {
			if val, ok := c[i][j]; ok {
				row = append(row, *val)
			} else {
				row = append(row, "")
			}
		}
		v = append(v, row)
	}
	return v
}

func parseColspanRowspan(n *html.Node) (int, int) {
	colspan := 1
	rowspan := 1
	for _, attr := range n.Attr {
		if strings.ToLower(attr.Key) == "colspan" {
			colspanRead, err := strconv.ParseInt(attr.Val, 10, 32)
			if err != nil {
				continue
			}
			colspan = int(colspanRead)
		}
		if strings.ToLower(attr.Key) == "rowspan" {
			rowspanRead, err := strconv.ParseInt(attr.Val, 10, 32)
			if err != nil {
				continue
			}
			rowspan = int(rowspanRead)
		}
	}
	return colspan, rowspan
}

func isFilled(c map[int]map[int](*string), i int, j int) bool {
	if c[i] == nil {
		c[i] = map[int]*string{}
	}
	if c[i][j] == nil {
		return false
	}
	return true
}

func editCell(c map[int]map[int]*string, i int, j int, str string) {
	if c[i] == nil {
		c[i] = map[int]*string{}
	}
	c[i][j] = &str
}

func editCellInRange(c map[int]map[int]*string, i int, j int, colspan int, rowspan int, str string) {
	for vi := 0; vi < rowspan; vi++ {
		for vj := 0; vj < colspan; vj++ {
			editCell(c, i+vi, j+vj, str)
		}
	}
}
