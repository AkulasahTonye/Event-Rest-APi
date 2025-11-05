.PHONY: run build test clean deps tidy migrate-up migrate-down

# Project variables
BINARY_NAME=api
DB_PATH=./data/sqlite.db
MIGRATIONS_DIR=./migrations

# Build the application
build:
	go build -o bin/$(BINARY_NAME) .

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build files
clean:
	go clean
	rm -f bin/$(BINARY_NAME)


# Tidy up dependencies
tidy:
	go mod tidy


