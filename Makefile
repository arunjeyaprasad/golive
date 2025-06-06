# Go parameters
BINARY_NAME=golive
GO=go

.PHONY: all build clean test coverage run

all: build

build:
	$(GO) build -o $(BINARY_NAME) .

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test:
	$(GO) test ./... -v

coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out

run:
	$(GO) run .

# Development targets
dev:
	$(GO) run .

lint:
	golangci-lint run

dockerbuild:
	docker build -t $(BINARY_NAME) .

dockerrun:
	docker run --rm -it -p 9090:9090 $(BINARY_NAME) --name $(BINARY_NAME)

dockerstop:
	docker stop $(BINARY_NAME)
	
.DEFAULT_GOAL := build