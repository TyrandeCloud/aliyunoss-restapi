.PHONY: build mod

mod:
	go mod download
	go mod tidy

run:
	go run ./test/main.go -listenPort=8085