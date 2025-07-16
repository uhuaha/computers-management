# Makefile for github.com/uhuaha/computers-management

## Run tests
test:
	go test ./...

## Run linting
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "▶ Running golangci-lint..."; \
		golangci-lint run ./...; \
	elif command -v staticcheck >/dev/null 2>&1; then \
		echo "▶ golangci-lint not found, falling back to staticcheck..."; \
		staticcheck ./...; \
	else \
		echo "⚠️ Neither golangci-lint nor staticcheck found; skipping lint"; \
	fi
