ENV ?= "development"
CONFIG_PATH ?= "./config"
DB_CONN ?= "postgres://postgres:123123@localhost/novelism?sslmode=disable"

run.local:
	ENV=${ENV} CONFIG_PATH=${CONFIG_PATH} go run ./cmd/api/

build:
	go build -o ./bin/ ./cmd/api/

build.docker:
	docker build --tag novelism-backend:${ENV} .

db.shell:
	docker exec -it postgresql_db psql -U postgres

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
