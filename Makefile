.PHONY: run version testdb test
.SILENT:
.ONESHELL:

PORT:=5000

VERSION := $(shell git describe --tags)

run: mydevto
	clear
	# source secret environment variables (passwords and stuff)
	export $$(grep -v '\(^$$\|^#\)' prod.env | xargs)
	export PORT=$(PORT)
	./mydevto

mydevto: $(shell find . -type f -name "*.go")
	go build -ldflags="-X 'github.com/math2001/mydevto/version.V=$(VERSION)'"

version:
	echo $(VERSION)

test:
	export $$(grep -v '\(^$$\|^#\)' test.env | xargs)
	go test -ldflags="-X 'github.com/math2001/mydevto/version.V=$(VERSION)-test'" ./...

testdb:
	export $$(grep -v '\(^$$\|^#\)' test.env | xargs)
	createdb $$DBNAME
	go run cmd/db/maketestdb.go
