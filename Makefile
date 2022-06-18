CURRENT_DIR = $(shell pwd)
LOCAL_BIN=$(CURRENT_DIR)/bin

bin-depth:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.1

run:
	@go run ./cmd/main.go

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --timeout 60s

swag:
	swag init -g cmd/main.go