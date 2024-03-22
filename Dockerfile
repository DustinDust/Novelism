
FROM  golang:1.21.1-alpine as build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o  ./bin/ ./cmd/api

FROM alpine:latest as build-release-stage
ENV ENV="development"
WORKDIR /
COPY --from=build-stage '/app/bin/api'  '/api'
COPY ./config/ ./config/
EXPOSE 8081

ENTRYPOINT [ "/api" ]
