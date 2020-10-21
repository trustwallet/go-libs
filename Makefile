GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

generate-coins:
	@echo "  >  Generating coin file"
	GOBIN=$(GOBIN) go run -tags=coins coin/gen.go
