.PHONY: install generate build run clean test dev help

install:
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest

generate:
	templ generate
	go run main.go build.go --generate

build: generate
	go build -o campaign

run: generate
	go run main.go build.go

dev: generate
	go run main.go build.go

test:
	@echo "Manual test:"
	@echo "1. go run main.go --generate"
	@echo "2. curl http://localhost:3000"
	@echo "3. curl -X POST http://localhost:3000/api/donations -d 'name=John&email=john@example.com&amount=50' -H 'Content-Type: application/x-www-form-urlencoded'"
	@echo "4. curl http://localhost:3000/api/stats"
	@echo "5. curl http://localhost:3000/api/recent-donors"

clean:
	rm -rf public/
	rm -f campaign
	rm -f templates/*_templ.go

.DEFAULT_GOAL := help
help:
	@echo "Available targets:"
	@echo "  make install   - Install dependencies"
	@echo "  make generate  - Generate templ code and static HTML"
	@echo "  make build     - Build binary"
	@echo "  make run       - Run server with static generation"
	@echo "  make dev       - Alias for run"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make test      - Show manual test commands"
