CURRENT_DIR = $(shell pwd)
LOCAL_BIN=$(CURRENT_DIR)/bin
VALUES=$(CURRENT_DIR)/.o3/k8s/values_local.yaml

ifndef PG_DSN
$(eval PG_DSN=$(shell cat $(VALUES) | grep -i "pg-dsn" -A1 | sed -n '2p;2q' | sed -e 's/[ \t]*value://g'))
endif

bin-depth:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.5.3
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.1

db\:create:
	$(LOCAL_BIN)/goose -dir migrations create "$(NAME)" sql

db\:up:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" up

db\:down:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" down

run:
	@go run ./cmd/main.go

test:
	PG_DSN=$(PG_DSN) VALUES=$(VALUES) go test -v ./...

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --timeout 60s