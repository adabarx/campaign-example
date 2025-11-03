#!/bin/bash

set -e

function install() {
    echo "Installing dependencies..."
    go mod download
    go install github.com/a-h/templ/cmd/templ@latest
}

function generate() {
    echo "Generating templ code and static HTML..."
    templ generate
    go run main.go build.go --generate
}

function build() {
    generate
    echo "Building binary..."
    go build -o campaign
}

function run() {
    generate
    echo "Running server..."
    go run main.go build.go
}

function dev() {
    run
}

function test() {
    echo "Manual test:"
    echo "1. go run main.go --generate"
    echo "2. curl http://localhost:3000"
    echo "3. curl -X POST http://localhost:3000/api/donations -d 'name=John&email=john@example.com&amount=50' -H 'Content-Type: application/x-www-form-urlencoded'"
    echo "4. curl http://localhost:3000/api/stats"
    echo "5. curl http://localhost:3000/api/recent-donors"
}

function clean() {
    echo "Cleaning build artifacts..."
    rm -rf public/
    rm -f campaign
    rm -f templates/*_templ.go
}

function help() {
    echo "Available commands:"
    echo "  ./build.sh install   - Install dependencies"
    echo "  ./build.sh generate  - Generate templ code and static HTML"
    echo "  ./build.sh build     - Build binary"
    echo "  ./build.sh run       - Run server with static generation"
    echo "  ./build.sh dev       - Alias for run"
    echo "  ./build.sh clean     - Clean build artifacts"
    echo "  ./build.sh test      - Show manual test commands"
}

case "$1" in
    install)
        install
        ;;
    generate)
        generate
        ;;
    build)
        build
        ;;
    run)
        run
        ;;
    dev)
        dev
        ;;
    test)
        test
        ;;
    clean)
        clean
        ;;
    help|"")
        help
        ;;
    *)
        echo "Unknown command: $1"
        help
        exit 1
        ;;
esac
