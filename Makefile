lint:
	golangci-lint run --config .golangci.yml
	@echo "Go lint passed successfully"

test:
	go test ./...
	@echo "All tests passed successfully"

.PHONY: lint test
