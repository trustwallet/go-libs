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
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v1.50.1
endif

go-lint:
	@echo "  >  Running golint"
	golangci-lint run ./...
