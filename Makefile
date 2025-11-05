.PHONY: test

ARGS ?= analyse -inputFile examples/decision.yaml
build:
	go build -o bin/ctac ./cmd/ctac

run: build
	./bin/ctac $(ARGS)

test:
	go test ./...

lint:
	go vet ./...
