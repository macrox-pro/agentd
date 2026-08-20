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
	@go build -o $(BIN) .

test:
	@go test ./... -race -count=1

test-short:
	@go test ./... -race -count=1 -short

e2e:
	@bash scripts/e2e-m1.sh
	@bash scripts/e2e-m2.sh
	@bash scripts/e2e-m3.sh
	@bash scripts/e2e-m4.sh
	@bash scripts/e2e-m5.sh
	@bash scripts/e2e-m6.sh

start: build
	@./$(BIN) daemon start

stop:
	@./$(BIN) daemon stop
