# rotato - simple spinner library for Go.
# See LICENSE file for copyright and license details.

FN		?= .

full: test lint

# Run tests
test:
	@echo '>> Testing $(BINARY_NAME)'
	@go test ./...
	@echo

# Run tests for a specific function
testfn:
	@echo '>> Testing function $(FN)'
	@go test -run $(FN) ./...

# Lint code with 'golangci-lint'
lint:
	@echo '>> Linting code'
	@go vet ./...
	golangci-lint run ./...
