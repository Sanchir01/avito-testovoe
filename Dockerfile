FROM golang:1.23-alpine AS builder

LABEL authors="sanchir"

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

RUN go version

COPY . .

RUN go mod download

RUN go build -o .bin/main cmd/main/main.go

ENTRYPOINT ["make", "run"]