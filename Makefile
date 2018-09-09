.PHONY: run
.SILENT:

PORT:=5000

run: mydevto
	env PORT=$(PORT) ./mydevto

mydevto: $(wildcard *.go) $(wildcard **/*.go)
	go build -i
