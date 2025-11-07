BIN := bin/ctac
ARGS ?= version
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.1.0-alpha.1)
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build run clean test lint

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/ctac

run: build
	./$(BIN) $(ARGS)

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/ctac

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BIN)
