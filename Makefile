
install-tools:
	@echo installing tools
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/matryer/moq@latest


build:
	@go build -o bin/ecom main.go

mock_generate:
	@go generate ./... 

test:
	@go test -v ./...

setup-local:
	@docker-compose up --wait  -d 

reset-local:
	@docker-compose down -v
	@docker-compose up --wait  -d 

run: build
	@./bin/ecom

run-local: setup-local
	go run main.go

migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a migration name using 'make migration-create name=<migration_name>'"; \
		exit 1; \
	fi
	goose -s -dir ./migrations create $(name) sql

migration-up:
	@go run migrations/migration.go up

migration-down:
	@go run migrations/migration.go down