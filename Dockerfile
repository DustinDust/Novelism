
FROM  golang:1.21.1-alpine as build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o  ./bin/ ./cmd/api

FROM alpine:latest as build-release-stage
ENV ENV="development"
ENV CONFIG_PATH="./config"
WORKDIR /
COPY --from=build-stage '/app/bin/api'  '/api'
COPY ./config/ "${CONFIG_PATH}/."
EXPOSE 8081

ENTRYPOINT [ "/api" ]
