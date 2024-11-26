.PHONY:
.SILENT:

run:
	chcp 65001
	go run ./cmd/app/main.go > temp.txt
