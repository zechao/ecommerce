

build:
	@go build -o bin/ecom main.go


test:
	@go test -v ./...

set-up:
	@docker-compose up --wait  -d 

run: build
	@./bin/ecom

run-local: set-up
	go run main.go
