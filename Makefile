PHONY:

SILENT:

MIGRATION_NAME ?= new_migration

PASSWORD ?= password
lint:
	golangci-lint run --config=.golangci.yaml

migrations-up:
	goose -dir migrations postgres "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable"  up

migrations-down:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable"  down


migrations-status:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable" status

migrations-new:
	goose -dir migrations create $(MIGRATION_NAME) sql


build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build	lint
	./.bin/main