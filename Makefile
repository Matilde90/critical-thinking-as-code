.PHONY: test

build:
	go build -o bin/ctac ./cmd/ctac

run: build
	./bin/ctac

test:
	go test ./...

lint:
	go vet ./...
