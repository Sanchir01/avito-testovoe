PHONY:
SILENT:

MIGRATION_NAME ?= new_migration

PASSWORD ?= password
lint:
	golangci-lint run --config=.golangci.yaml

migrations-up:
	goose -dir migrations postgres "host=localhost user=postgres password=avito port=5432 dbname=postgres sslmode=disable"  up

migrations-down:
	goose -dir migrations postgres  "host=localhost user=postgres password=avito port=5432 dbname=postgres sslmode=disable"  down

migrations-status:
	goose -dir migrations postgres  "host=localhost user=postgres password=avito port=5432 dbname=postgres sslmode=disable" status

migrations-new:
	goose -dir migrations create $(MIGRATION_NAME) sql


build:
	go build -o ./.bin/main ./cmd/main/main.go

swag:
	swag init -g cmd/main/main.go

run: build	lint  swag
	./.bin/main

docker:
		docker build -t avito .
		docker-compose up --build app

testing:
		go test -v ./test/...

seed:
	go run cmd/seed/main.go