test: test-network
	go test -v ./...

test-network:
	cd ./network; \
	go test -v ./...; \

## golint: Run linter.
lint: go-lint-install go-lint

go-lint-install:
ifeq (,$(wildcard test -f bin/golangci-lint))
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

go-lint:
	@echo "  >  Running golint"
	bin/golangci-lint run --timeout=2m
