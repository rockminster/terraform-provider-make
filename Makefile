default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build the provider
.PHONY: build
build:
	go build -v .

# Install the provider locally for testing
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rockminster/make/0.1.0/linux_amd64
	mv terraform-provider-make ~/.terraform.d/plugins/registry.terraform.io/rockminster/make/0.1.0/linux_amd64/

# Generate documentation
.PHONY: docs
docs:
	go generate ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-make

# Lint the code
.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Installing..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Test the code
.PHONY: test
test:
	go test ./... -v

# Run all checks
.PHONY: check
check: fmt lint test

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build     - Build the provider binary"
	@echo "  install   - Install the provider locally for testing"
	@echo "  test      - Run unit tests"
	@echo "  testacc   - Run acceptance tests"
	@echo "  docs      - Generate documentation"
	@echo "  lint      - Run linters"
	@echo "  fmt       - Format code"
	@echo "  check     - Run all checks (fmt, lint, test)"
	@echo "  clean     - Clean build artifacts"
	@echo "  help      - Show this help message"