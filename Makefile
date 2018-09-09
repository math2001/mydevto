.PHONY: run
.SILENT:

run: mydevto
	./mydevto

mydevto: main.go
	go build -i
