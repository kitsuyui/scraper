# Scraper: Swiss Army Knife for Web scraping

[![CircleCI Status](https://circleci.com/gh/kitsuyui/scraper.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/kitsuyui/scraper)
[![Codecov Status](https://img.shields.io/codecov/c/github/kitsuyui/scraper.svg)](https://codecov.io/github/kitsuyui/scraper/)
[![Github All Releases](https://img.shields.io/github/downloads/kitsuyui/scraper/total.svg)](https://github.com/kitsuyui/scraper/releases/latest)

# Installation

You can download executable static binary from https://github.com/kitsuyui/scraper/releases/latest .

- typically Linux ... download scraper_darwin_amd64.
- typically macOS ... download scraper_linux_amd64.
- typically Windows ... download scraper_windows_amd64.exe

and then just you have to do is adding PATH and permitting it to executable.

```
$ target=scraper_darwin_amd64  # executable binary you downloaded.
$ cp "./$target" /usr/bin/scraper
$ chmod +x /usr/bin/scraper
```

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

### Regular Expression Scraping

```
$ wget http://example.com/ -O index.html
$ cat <<'EOF'> scraper-config.json
[
  {"type": "regex", "label": "Description", "query": "This domain is .+? to be used for .+?\\."}
]
EOF
$ scraper < index.html
[
   {
     "type": "regex",
     "label": "Description",
     "query": "This domain is .+? to be used for .+?\\.",
     "results": [
       "This domain is established to be used for illustrative examples in documents."
     ]
   }
]
```

### Composable

Obviously these recipes are mixable.
You can scrape all at once with one-config by writing these target.

### Server mode

### startup

```console
$ scraper server
```

### Put configuration

```console
$ curl -X PUT localhost:8080/example.com -d@- <<'EOT'
[{"type": "xpath", "label": "LinkURL", "query": "//a/@href"}]
EOT
```

## Show configuration

```console
$ curl -X GET localhost:8080/example.com
[{"type": "xpath", "label": "LinkURL", "query": "//a/@href"}]
```

## Scraping

```console
$ curl example.com | curl -X POST localhost:8080/example.com -d @-
[
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
$ ./bin/build.sh
```

#### with Docker

```console
$ docker run --rm -v "$(pwd)":/root -w /root tcnksm/gox sh -c "./bin/build.sh"
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
