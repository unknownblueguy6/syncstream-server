.PHONY: tools
tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.57.2

.PHONY: lint
lint:
	golangci-lint run --timeout 90m ./...

.PHONY: format
format:
	golangci-lint run --fix ./...

.PHONY: test
test:
	go test -cover ./...
