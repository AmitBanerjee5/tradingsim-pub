# Makefile for tradingsim

BINARY=tradingsim
SRC=main.go
DOCKER_IMAGE=tradingsim:latest
CONFIG=configprocessor/sampleconfig.json

.PHONY: all build run clean fmt lint docker-run docker-clean docker-exec

all: build

build: clean
	go build -o $(BINARY) $(SRC)
	chmod +x $(BINARY)

run: build
	cat $(CONFIG)|./$(BINARY)

fmt:
	go fmt ./...

lint:
	golint ./...

clean:
	rm -f $(BINARY)

docker-run:
	docker/createdocker.sh

docker-clean:
	docker/removedocker.sh

docker-exec:
	docker/execdocker.sh

# Example usage:
# make build          # Build the binary
# make run            # Build and run the binary
# make fmt            # Format code
# make lint           # Lint code
# make clean          # Remove binary
# make docker-run     # Build and run Docker container
# make docker-clean   # Stop and remove Docker container and image
# make docker-exec    # Exec into running Docker container
