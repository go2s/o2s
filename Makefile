lint:
	golangci-lint run

format:
		goimports -w -l .
		go fmt ./...

test:
		go test ./... -v

all: format lint test