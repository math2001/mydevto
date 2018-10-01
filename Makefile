.PHONY: run version testdb
.SILENT:
.ONESHELL:

PORT:=5000

VERSION := $(shell git describe --tags)

run: mydevto
	echo $(VERSION)
	clear
	# source secret environment variables (passwords and stuff)
	export $$(grep -v '\(^$$\|^#\)' prod.env | xargs)
	export PORT=$(PORT)
	./mydevto

mydevto: $(shell find . -type f -name "*.go")
	go build -v -ldflags="-X main.version=$(VERSION)"

version:
	echo $(VERSION)

testdb:
	export $$(grep -v '\(^$$\|^#\)' test.env | xargs)
	createdb $$DBNAME || echo "-> Error ignored. Creating schema..."
	psql -U $$DBLOGIN -d $$DBNAME -f ./bin/schema.pgsql -f ./bin/populate_test.pgsql
