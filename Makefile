test: test-network
	go test -v ./...

test-network:
	cd ./network; \
	go test -v ./...; \

## golint: Run linter.
lint: go-lint-install go-lint

go-lint-install:
ifeq (,$(shell which golangci-lint))
	@echo "  >  Installing golint"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0
endif

go-lint:
	@echo "  >  Running golint"
	golangci-lint run ./...
