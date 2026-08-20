.PHONY: generate generate-grpc lint test test-short e2e docs-check build start stop

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

# EN/RU user-doc filename parity (same basenames under docs/en and docs/ru).
docs-check:
	@tmp_en=$$(mktemp) tmp_ru=$$(mktemp); \
	trap 'rm -f "$$tmp_en" "$$tmp_ru"' EXIT; \
	(cd docs/en && ls *.md | sort) >"$$tmp_en"; \
	(cd docs/ru && ls *.md | sort) >"$$tmp_ru"; \
	diff -u "$$tmp_en" "$$tmp_ru" \
		|| { echo "docs-check: docs/en and docs/ru must have the same *.md basenames" >&2; exit 1; }; \
	echo "docs-check: ok"

start: build
	@./$(BIN) daemon start

stop:
	@./$(BIN) daemon stop
