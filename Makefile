.PHONY: lint

lint:
	golangci-lint run

test:
	go test ./...

install:
	go build -o git-auto-sync .
	cd daemon && go build -o git-auto-sync-daemon .
	sudo mv git-auto-sync ${GOPATH}
	cd daemon && sudo mv git-auto-sync-daemon ${GOPATH}
