.PHONY: build
build:
	go build -v ./cmd/medos && migrate -path migrations -database "postgres://postgres:qwe@localhost/go_medos?sslmode=disable" up
.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build