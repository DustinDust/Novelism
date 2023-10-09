ENV ?= "development"

run.dev:
	ENV=${ENV} go run ./cmd/api/ 

build:
	go build -o ./bin/ ./cmd/api/

start:
	make build && ./bin/api