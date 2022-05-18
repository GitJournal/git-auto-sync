.PHONY: lint
	
lint:
	golangci-lint run

test:
	go test ./...

install:
	go install .
	cd daemon && go build -o git-auto-sync-daemon .
	cd daemon && mv git-auto-sync-daemon ${GOPATH}/bin/
