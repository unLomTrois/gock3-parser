run:
	go run ./cmd/main.go

build:
	cd ./cmd && go build main.go

test:
	go test ./...

.DEFAULT_GOAL := build
