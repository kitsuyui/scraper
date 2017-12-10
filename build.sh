#!/bin/sh
go get -d ./...
go get -d github.com/inconshreveable/mousetrap
go get -d github.com/Microsoft/go-winio
gox -ldflags '-w -s' -output='build/scraper_{{.OS}}_{{.Arch}}'
