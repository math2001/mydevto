.PHONY: run version testdb test
.SILENT:
.ONESHELL:

PORT:=5000

VERSION := $(shell git describe --tags)

LDFLAGS := -ldflags="-X 'github.com/math2001/mydevto/services/buildinfos.V=$(VERSION)$$VERSION_EXTRA'"

run: mydevto
	clear
	# source secret environment variables (passwords and stuff)
	export $$(grep -v '\(^$$\|^#\)' prod.env | xargs)
	export PORT=$(PORT)
	./mydevto

mydevto: $(shell find . -type f -name "*.go")
	go build $(LDFLAGS)

version:
	echo $(VERSION)

test:
	export $$(grep -v '\(^$$\|^#\)' test.env | xargs)
	export VERSION_EXTRA='-testing'
	go test $(LDFLAGS) ./...

testdb:
	export $$(grep -v '\(^$$\|^#\)' test.env | xargs)
	createdb $$DBNAME
	go run cmd/db/maketestdb.go
