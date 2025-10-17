# ========================
# 🧱 VARIÁVEIS
# ========================
APP_NAME=auction
IMAGE_NAME=$(APP_NAME):latest
PORT=8080
PKG := $(shell go list ./...)
BIN=server
MAIN=./cmd/auction/main.go
COMPOSE_FILE=./docker-compose.yml
MOCKERY_VERSION := v2.53.2
BIN_MOCKERY := $(shell go env GOPATH)/bin/mockery

# Coverage files
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
COVERAGE_UNIT_FILE := coverage_unit.out
COVERAGE_UNIT_HTML := coverage_unit.html
COVERAGE_INTEGRATION_FILE := coverage_integration.out
COVERAGE_INTEGRATION_HTML := coverage_integration.html

# ========================
# 🧩 AJUDA
# ========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "📦 Go commands:"
	@echo "  make start                         -> Run the application locally"
	@echo "  make build                         -> Build binary"
	@echo "  make test                          -> Run all tests (unit + integration)"
	@echo "  make test-unit                     -> Run unit tests only"
	@echo "  make test-integration              -> Run integration tests only"
	@echo "  make coverage                      -> Show combined coverage summary"
	@echo "  make coverage-html                 -> Generate combined HTML coverage report"
	@echo "  make coverage-unit                 -> Show unit test coverage summary"
	@echo "  make coverage-html-unit            -> Generate HTML report for unit tests"
	@echo "  make coverage-integration          -> Show integration test coverage summary"
	@echo "  make coverage-html-integration     -> Generate HTML report for integration tests"
	@echo "  make clear                         -> Remove binary and coverage files"
	@echo "  make fmt                           -> Format code"
	@echo "  make lint                          -> Run linter"
	@echo ""
	@echo "🐳 Docker commands:"
	@echo "  make docker-build                  -> Build Docker image"
	@echo "  make docker-run                    -> Run Docker image locally"
	@echo "  make up                            -> Run service via docker-compose"
	@echo "  make down                          -> Stop docker-compose"
	@echo "  make logs                          -> Follow logs"
	@echo "  make status                        -> Check container status"
	@echo "  make clean                         -> Remove containers, volumes and orphans"
	@echo ""
	@echo "🧪 Mockery:"
	@echo "  make install-mockery               -> Install mockery"
	@echo "  make reinstall-mockery             -> Force reinstall mockery"
	@echo "  make generate-mocks                -> Generate mocks via go:generate"

# ========================
# 🚀 GO COMMANDS
# ========================
.PHONY: start build test test-unit test-integration coverage coverage-html coverage-unit coverage-html-unit coverage-integration coverage-html-integration clear fmt lint

start:
	@echo "🚀 Starting $(APP_NAME)..."
	go run $(MAIN)

build:
	@echo "🔨 Building binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) $(MAIN)

# Run all tests (unit + integration)
test:
	@echo "🧪 Running all tests (unit + integration)..."
	go test -v -coverpkg=./... -cover -coverprofile=$(COVERAGE_FILE) $(PKG)

# Unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	go test -v -coverpkg=./... -cover -tags=unit -coverprofile=$(COVERAGE_UNIT_FILE) $(PKG)

# Integration tests only
test-integration:
	@echo "🧪 Running integration tests..."
	go test -v -coverpkg=./... -cover -tags=integration -coverprofile=$(COVERAGE_INTEGRATION_FILE) $(PKG)

# Coverage reports - combined
coverage:
	@echo "📊 Combined coverage summary:"
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: test
	@echo "🌐 Generating combined HTML coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@open $(COVERAGE_HTML) || echo "To open manually: $(COVERAGE_HTML)"

# Coverage reports - unit only
coverage-unit: test-unit
	@echo "📊 Unit test coverage summary:"
	go tool cover -func=$(COVERAGE_UNIT_FILE)

coverage-html-unit: test-unit
	@echo "🌐 Generating HTML report for unit tests..."
	go tool cover -html=$(COVERAGE_UNIT_FILE) -o $(COVERAGE_UNIT_HTML)
	@open $(COVERAGE_UNIT_HTML) || echo "To open manually: $(COVERAGE_UNIT_HTML)"

# Coverage reports - integration only
coverage-integration: test-integration
	@echo "📊 Integration test coverage summary:"
	go tool cover -func=$(COVERAGE_INTEGRATION_FILE)

coverage-html-integration: test-integration
	@echo "🌐 Generating HTML report for integration tests..."
	go tool cover -html=$(COVERAGE_INTEGRATION_FILE) -o $(COVERAGE_INTEGRATION_HTML)
	@open $(COVERAGE_INTEGRATION_HTML) || echo "To open manually: $(COVERAGE_INTEGRATION_HTML)"

clear:
	@echo "🧹 Cleaning binary and coverage files..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML) $(COVERAGE_UNIT_FILE) $(COVERAGE_UNIT_HTML) $(COVERAGE_INTEGRATION_FILE) $(COVERAGE_INTEGRATION_HTML) $(BIN)

fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Running linter..."
	golangci-lint run || true

# ========================
# 🐳 DOCKER COMMANDS
# ========================
.PHONY: docker-build docker-run up down logs status stop clean
docker-build:
	@echo "🐳 Building Docker image $(IMAGE_NAME)..."
	docker build -t $(IMAGE_NAME) -f ./Dockerfile .

docker-run:
	@echo "🚀 Running container on port $(PORT)..."
	docker run --rm -p $(PORT):$(PORT) --env-file .env $(IMAGE_NAME)

up:
	@echo "🚀 Starting $(APP_NAME) via docker-compose..."
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	@echo "🧹 Stopping docker-compose..."
	docker compose -f $(COMPOSE_FILE) down

stop:
	@echo "⏹️  Stopping containers..."
	docker compose -f $(COMPOSE_FILE) stop

status:
	docker compose -f $(COMPOSE_FILE) ps

logs:
	docker compose -f $(COMPOSE_FILE) logs --follow

clean:
	@echo "🧹 Removing containers, volumes, and orphans..."
	docker compose -f $(COMPOSE_FILE) down -v --remove-orphans

# ========================
# 🧪 MOCKERY
# ========================
.PHONY: install-mockery reinstall-mockery generate-mocks
install-mockery:
	@echo "📦 Installing mockery..."
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

reinstall-mockery:
	@echo "🔁 Reinstalling mockery..."
	rm -f $(BIN_MOCKERY)
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

generate-mocks:
	@echo "⚙️ Generating mocks..."
	go generate ./...
