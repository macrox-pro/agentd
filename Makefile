.PHONY: generate generate-grpc lint test test-short e2e build start stop

BIN ?= agentd

generate: generate-grpc

generate-grpc:
	@buf lint api/
	@buf generate

lint:
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || go vet ./...
	@buf lint api/

build:
	@go build -ldflags "-X github.com/macrox-pro/agentd/internal/version.Version=dev" -o $(BIN) .

test:
	@go test ./... -race -count=1

test-short:
	@go test ./... -race -count=1 -short

e2e:
	@for s in $$(ls scripts/e2e-m*.sh | sort -V); do bash "$$s" || exit 1; done

start: build
	@./$(BIN) daemon start

stop:
	@./$(BIN) daemon stop
