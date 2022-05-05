.PHONY: lint
	
lint: export CGO_ENABLED = 0
lint:
	golangci-lint run

test:
	go test ./...
