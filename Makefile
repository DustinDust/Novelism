ENV ?= "development"
DB_CONN ?= "postgres://postgres:123123@localhost/novelism?sslmode=disable"

run.dev:
	ENV=${ENV} go run ./cmd/api/

build:
	go build -o ./bin/ ./cmd/api/

start:
	make build && ./bin/api

migrate.create:
	migrate create -seq -ext .sql -dir ./migrations ${NAME}

migrate.up:
	migrate -path ./migrations -database ${DB_CONN} up ${STEP}

migrate.down:
	migrate -path ./migrations -database ${DB_CONN} down ${STEP}

migrate.force:
	migrate -path ./migrations -database ${DB_CONN} force ${VERSION}

migrate.drop:
	migrate -path ./migrations -database ${DB_CONN} drop

migrate.version:
	migrate -path ./migrations -database ${DB_CONN} version
