FROM golang:1.24.0-alpine AS builder

WORKDIR /app

RUN apk --no-cache add make

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.4

RUN go version

COPY . .

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN go build -o .bin/main cmd/main/main.go

ENTRYPOINT ["make", "run"]