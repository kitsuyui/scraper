#!/bin/sh
go get
CGO_ENABLE=0 gox -ldflags '-w -s' -output='build/scraper_{{.OS}}_{{.Arch}}'
