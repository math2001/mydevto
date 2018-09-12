.PHONY: run
.SILENT:
.ONESHELL:

PORT:=5000

run: mydevto
	clear
	# source secret environment variables (passwords and stuff)
	export $$(grep -v '\(^$$\|^#\)' .env | xargs)
	export PORT=$(PORT)
	./mydevto

mydevto: $(wildcard *.go) $(wildcard **/*.go)
	go build -i
