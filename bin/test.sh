#!/bin/sh
go get
go get -d github.com/inconshreveable/mousetrap
go get -d github.com/Microsoft/go-winio
go test -v ./...
