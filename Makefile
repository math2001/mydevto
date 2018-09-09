.PHONY: run
.SILENT:

PORT:=5000

run: mydevto
	env PORT=$(PORT) ./mydevto

mydevto: main.go
	go build -i
