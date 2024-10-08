.PHONY: build
build:
	go build -v ./cmd/medods && migrate -path migrations -database "postgres://postgres:qwe@localhost/go_medods?sslmode=disable" up
.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build