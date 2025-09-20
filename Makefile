include deploy/.env
LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	@if [ ! -f "$(LOCAL_BIN)/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
	else \
		echo "golangci-lint already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/goose" ]; then \
		echo "Installing goose..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.0; \
	else \
		echo "goose already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/minimock" ]; then \
		echo "Installing minimock..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.4.5 ; \
	else \
		echo "minimock already installed."; \
	fi
	@if [ ! -f "$(LOCAL_BIN)/mockgen" ]; then \
		echo "Installing mockgen..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1; \
	else \
		echo "mockgen already installed."; \
	fi


lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


local migration-create:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create $(name) sql sql

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v


local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v


local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go
copy-to-server:
	scp service_linux root@$(IP_SERVER):

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t $(REGESTRY)/server:v0.0.1 -f deploy/Dockerfile .
	docker login -u $(USERNAME) -p $(PASSWORD) $(REGESTRY)
	docker push $(REGESTRY)/server:v0.0.1

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/Ippolid/auth/internal/service/...,github.com/Ippolid/auth/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Ippolid/auth/internal/service/...,github.com/Ippolid/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore




