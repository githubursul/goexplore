APP_NAME := goexplore

.PHONY: build build-static run run-static

build:
	go build -o $(APP_NAME) ./...

build-static:
	go build -tags=static -o $(APP_NAME) ./...

run:
	go run .

run-static:
	go run -tags=static .
