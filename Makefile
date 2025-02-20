.PHONY: build run test clean mod-download mod-tidy mod-verify mod-update new-command

# Binary name
BINARY_NAME=crosslab

# Build the application using goreleaser
build:
	goreleaser build --snapshot --clean --single-target

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f bin/${BINARY_NAME}

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Run the application in development mode
dev:
	go build -o ./bin/${BINARY_NAME}
	./bin/${BINARY_NAME}

# Create a new release using goreleaser
release:
	goreleaser release --rm-dist

# Install goreleaser (you may want to run this first)
install-goreleaser:
	go install github.com/goreleaser/goreleaser@latest

# Download all dependencies
mod-download:
	go mod download

# Tidy the go.mod and go.sum files
mod-tidy:
	go mod tidy

# Verify dependencies
mod-verify:
	go mod verify

# Update all dependencies
mod-update:
	go get -u ./...
	go mod tidy

# Initialize a new module
mod-init:
	go mod init ${BINARY_NAME}

# Show module dependencies
mod-graph:
	go mod graph

install-binary:
	- gh release download --pattern "crosslocal-*" --dir bin

