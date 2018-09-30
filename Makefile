.PHONY: run
.SILENT:
.ONESHELL:

PORT:=5000

VERSION := $(shell git describe --tags)

run: mydevto
	echo $(VERSION)
	clear
	# source secret environment variables (passwords and stuff)
	export $$(grep -v '\(^$$\|^#\)' .env | xargs)
	export PORT=$(PORT)
	./mydevto

mydevto: $(shell find . -type f -name "*.go")
	go build -v -ldflags="-X main.version=$(VERSION)"
