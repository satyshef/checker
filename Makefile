APP_PATH="checker"

.PHONY: build
build:
	go build -modfile go.mod -v -o ${APP_PATH} ./cmd/app  

.DEFAULT_GOAL := build
