GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

test: test-networks
	go test -v ./...

test-networks:
	cd ./networks; \
	go test -v ./...; \

generate-coins:
	@echo "  >  Generating coin file"
	GOBIN=$(GOBIN) go run -tags=coins coin/gen.go
