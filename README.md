# Scraper: Swiss Army Knife for Web scraping

[![CircleCI Status](https://circleci.com/gh/kitsuyui/scraper.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/kitsuyui/scraper)
[![Codecov Status](https://codecov.io/gh/kitsuyu/scraper/branch/master/graph/badge.svg)](https://codecov.io/github/kitsuyui/scraper/)

## Usage

### Scraping with XPath / CSS Selector

```
$ wget http://example.com/ -O index.html
$ cat <<'EOF'> scraper-config.json
[
  {"type": "css", "label": "Title", "query": "h1"},
  {"type": "xpath", "label": "LinkURL", "query": "//a/@href"}
]
EOF
$ scraper < index.html
[
   {
     "type": "css",
     "label": "Title",
     "query": "h1",
     "results": [
       "Example Domain"
     ]
   },
   {
     "type": "xpath",
     "label": "LinkURL",
     "query": "//a/@href",
     "results": [
       "http://www.iana.org/domains/example"
     ]
   }
 ]
```

### Table Scraping

Table element scraping is also available.
Obviously it supports colspan/rowspan attributes.

```
$ cat <<'EOF'> table.html
<table border=1>
  <tr><td>a</td><td>b</td><td>c</td><td>d</td></tr>
  <tr><td>e</td><td rowspan="2" colspan="2">f</td><td>g</td></tr>
  <tr><td>h</td><td>i</td></tr>
</table>
EOF
$ cat <<'EOF'> table-scraper-config.json
[
  {"type": "table-xpath", "label": "Tables", "query": "//table"}
]
EOF
$ scraper -c table-scraper-config.json < table.html
[
   {
     "type": "table-xpath",
     "label": "Tables",
     "query": "//table",
     "results": [
       [
         ["a", "b", "c", "d"],
         ["e", "f", "f", "g"],
         ["h", "f", "f", "i"]
       ]
     ]
   }
 ]
```

### Options

```
$ scraper -h
Scraper: Swiss Army Knife for Web scraping

Usage:
  scraper [flags]
  scraper [command]

Available Commands:
  help        Help about any command
  validate    Validate config file

Flags:
  -c, --configFile string   config file (default "scraper-config.json")
  -h, --help                help for scraper
  -i, --in string           input file (default STDIN)
  -o, --out string          output file (default STDOUT)

Use "scraper [command] --help" for more information about a command.
```

## Build

### For developping

```
$ go get -d ./src/...
$ go build -o build/scraper main.go
```

### For cross-platform

```
$ ./build.sh
```

#### with Docker

```console
$ docker run --rm -v "$(pwd)":/root -w /root tcnksm/gox sh -c "./build.sh"
```

## LICENSE

### Source

The 3-Clause BSD License. See also LISENCE file.

### statically linked libraries:

- [golang/go](https://github.com/golang/go/) ... [BSD 3-clause "New" or "Revised" License](https://github.com/golang/go/blob/master/LICENSE)
-	[antchfx/xpath](https://github.com/antchfx/xpath/) ... [MIT License](https://github.com/antchfx/xpath/blob/master/LICENSE)
-	[antchfx/xquery](https://github.com/antchfx/xquery/) ... [MIT License](https://github.com/antchfx/xquery/blob/master/LICENSE)
-	[andybalholm/cascadia](https://github.com/andybalholm/cascadia/) ... [BSD 2-clause "Simplified" License](https://github.com/andybalholm/cascadia/blob/master/LICENSE)
- [spf13/cobra](https://github.com/spf13/cobra/) ... [Apache License 2.0](https://github.com/spf13/cobra/blob/master/LICENSE.txt)
  - (windows) [inconshreveable/mousetrap](https://github.com/inconshreveable/mousetrap/) ... [Apache License 2.0](https://github.com/inconshreveable/mousetrap/blob/master/LICENSE)
  - (windows) [Microsoft/go-winio](https://github.com/Microsoft/go-winio/) ... [MIT License](https://github.com/Microsoft/go-winio/blob/master/LICENSE)
