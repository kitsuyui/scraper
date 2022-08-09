# Scraper: Swiss Army Knife for Web scraping

[![Codecov Status](https://img.shields.io/codecov/c/github/kitsuyui/scraper.svg)](https://codecov.io/github/kitsuyui/scraper/)
[![Github All Releases](https://img.shields.io/github/downloads/kitsuyui/scraper/total.svg)](https://github.com/kitsuyui/scraper/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/kitsuyui/scraper.svg)](https://hub.docker.com/r/kitsuyui/scraper/)

# Installation

## Download binary

You can download executable static binary from https://github.com/kitsuyui/scraper/releases/latest .

- typically Linux ... download scraper_linux_amd64.
- typically macOS ... download scraper_darwin_amd64.
- typically Windows ... download scraper_windows_amd64.exe

and then just you have to do is adding PATH and permitting it to executable.

```
$ target=scraper_darwin_amd64  # executable binary you downloaded.
$ cp "./$target" /usr/bin/scraper
$ chmod +x /usr/bin/scraper
```

## Install as go module

```console
$ go install github.com/kitsuyui/scraper@latest
```

# Usage

## Scraping with XPath / CSS Selector

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

## Table Scraping

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

## Regular Expression Scraping

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

## Composable

Obviously these recipes are mixable.
You can scrape all at once with one-config by writing these target.

## Validate configuration

Configurations can be checked by using subcommand `validate`.
This sub-command also checks format of each types.
For example,

- `regex` type will be checked the query is valid regular expression.
- `xpath` type will be checked the query is valid XPath.

```
$ scraper validate
Error: Recipe[0]: expected identifier, found < instead
exit status 1
```

## Server mode

### startup

```console
$ scraper server
```

### HTTP PUT: Define/Update Configuration

```console
$ curl -X PUT localhost:8080/example.com -d@- <<'EOT'
[{"type": "xpath", "label": "LinkURL", "query": "//a/@href"}]
EOT
```

### HTTP GET: Show configuration

```console
$ curl -X GET localhost:8080/example.com
[{"type": "xpath", "label": "LinkURL", "query": "//a/@href"}]
```

### HTTP POST: Scraping

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

## Docker

Docker image for server mode: [kitsuyui/scraper](https://hub.docker.com/r/kitsuyui/scraper/)

```
$ docker run -p 8080:8080 -it kitsuyui/scraper
```

### Options

```
$ scraper -h
scraper: Swiss army knife for web scraping

Usage:
 scraper [-c=<conf>] [-i=<input>]
 scraper validate [-c=<conf>]
 scraper server [-d=<conf-dir>] [-H=<host>] [-p=<port>]

Options:
 -c=<conf> --config=<conf>            Configuration file [default: scraper-config.json].
 -i=<input> --input=<input>           Input file specified instead of STDIN.
 -o=<output> --output=<output>        Output file specified instead of STDOUT.
 -H=<host> --host=<host>              Server mode host [default: 127.0.0.1].
 -p=<port> --port=<port>              Server mode port [default: 8080].
 -d=<conf-dir> --conf-dir=<conf-dir>  Configuration directory for server mode [default: .].
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

The 3-Clause BSD License. See also LICENSE file.

### statically linked libraries:

- [golang/go](https://github.com/golang/go/) ... [BSD 3-clause "New" or "Revised" License](https://github.com/golang/go/blob/master/LICENSE)
- [antchfx/xpath](https://github.com/antchfx/xpath/) ... [MIT License](https://github.com/antchfx/xpath/blob/master/LICENSE)
- [antchfx/xquery](https://github.com/antchfx/xquery/) ... [MIT License](https://github.com/antchfx/xquery/blob/master/LICENSE)
- [andybalholm/cascadia](https://github.com/andybalholm/cascadia/) ... [BSD 2-clause "Simplified" License](https://github.com/andybalholm/cascadia/blob/master/LICENSE)
- [docopt/docopt-go](https://github.com/docopt/docopt.go) ... [MIT License](https://github.com/docopt/docopt.go/blob/master/LICENSE)
