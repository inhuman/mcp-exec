.PHONY: build test test-short vet vulncheck vendor docker

build:
	go build -mod=vendor -trimpath -o mcp-exec ./cmd/mcp-exec

test:
	go test ./...

test-short:
	go test -short ./...

vet:
	go vet ./...

vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

vendor:
	go mod tidy && go mod vendor

docker:
	docker build -t idconstruct/mcp-exec:dev .
