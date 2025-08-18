# Makefile for github.com/uhuaha/computers-management

## Run unit tests
test:
	go test -cover ./...

## Run integration tests (require Docker + Build-Tag)
test-integration:
	go test -v -tags=integration ./internal/integration/...

## Run linting
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "▶ Running golangci-lint..."; \
		golangci-lint run --build-tags=integration ./...; \
	elif command -v staticcheck >/dev/null 2>&1; then \
		echo "▶ golangci-lint not found, falling back to staticcheck..."; \
		staticcheck ./...; \
	else \
		echo "⚠️ Neither golangci-lint nor staticcheck found; skipping lint"; \
	fi

## Run formatting
fmt:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "▶ Formatting files..."; \
		golangci-lint run --fix --build-tags=integration ./...; \
	else \
		echo "⚠️ No golangci-lint found; skipping formatting"; \
	fi
