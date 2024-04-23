ENV ?= "development"
DB_CONN ?= "postgres://postgres:123123@localhost/novelism?sslmode=disable"
TASK ?= main.go

run.local:
	ENV=${ENV} go run ./cmd/api/

build.local:
	go build -o ./bin/ ./cmd/api/

run.container:
	docker run -p 8081:8081 --name novelism-backend-${ENV} novelism-backend:${ENV}

build.container:
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

# e.g make run.task ENV="production" TASK=cmd/task/populate_user_status.go
run.task:
	ENV=${ENV} go run ${TASK}

# Define variable for filtering images
IMAGE_FILTER := "novelism"

# Target to clean images with 'novelism' in name
# Currently not working
clean.images:
  # Get a list of images matching the filter
  IMG := $$(docker images --format '{{.Repository}}:{{.Tag}}' | grep $(IMAGE_FILTER))
  echo $(IMG)
