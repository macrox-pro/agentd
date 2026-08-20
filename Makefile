.PHONY: generate generate-grpc lint test test-short

generate: generate-grpc

generate-grpc:
	@buf lint api/
	@buf generate

lint:
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || go vet ./...
	@buf lint api/

test:
	@go test ./... -race -count=1

test-short:
	@go test ./... -race -count=1 -short
