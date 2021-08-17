GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

test: test-network
	go test -v ./...

test-network:
	cd ./network; \
	go test -v ./...; \

generate-coins:
	@echo "  >  Generating coin file"
	GOBIN=$(GOBIN) go run -tags=coins coin/gen.go
	goimports -w coin/coins.go
