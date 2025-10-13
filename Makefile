# ========================
# 🧱 VARIÁVEIS
# ========================
APP_NAME=auction
IMAGE_NAME=$(APP_NAME):latest
PORT=8080
PKG := $(shell go list ./...)
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
BIN=server
MAIN=./cmd/auction/main.go
COMPOSE_FILE=./docker-compose.yml
MOCKERY_VERSION := v2.53.2
BIN_MOCKERY := $(shell go env GOPATH)/bin/mockery

# ========================
# 🧩 AJUDA
# ========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "📦 Go commands:"
	@echo "  make start                 -> Run the application locally"
	@echo "  make build                 -> Build binary"
	@echo "  make test                  -> Run tests with coverage"
	@echo "  make coverage              -> Show coverage summary"
	@echo "  make coverage-html         -> Generate HTML coverage report"
	@echo "  make clear                 -> Remove binary and coverage files"
	@echo "  make fmt                   -> Format code"
	@echo "  make lint                  -> Run linter"
	@echo ""
	@echo "🐳 Docker commands:"
	@echo "  make docker-build          -> Build Docker image"
	@echo "  make docker-run            -> Run Docker image locally"
	@echo "  make up                    -> Run service via docker-compose"
	@echo "  make down                  -> Stop docker-compose"
	@echo "  make logs                  -> Follow logs"
	@echo "  make status                -> Check container status"
	@echo "  make clean                 -> Remove containers, volumes and orphans"
	@echo ""
	@echo "🧪 Mockery:"
	@echo "  make install-mockery       -> Install mockery"
	@echo "  make reinstall-mockery     -> Force reinstall mockery"
	@echo "  make generate-mocks        -> Generate mocks via go:generate"

# ========================
# 🚀 GO COMMANDS
# ========================
.PHONY: start build test coverage coverage-html clear fmt lint
start:
	@echo "🚀 Starting $(APP_NAME)..."
	go run $(MAIN)

build:
	@echo "🔨 Building binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) $(MAIN)

test:
	@echo "🧪 Running tests..."
	go test -v -cover -coverprofile=$(COVERAGE_FILE) $(PKG)

coverage:
	@echo "📊 Test coverage summary:"
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: test
	@echo "🌐 Generating HTML coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@open $(COVERAGE_HTML) || echo "To open manually: $(COVERAGE_HTML)"

clear:
	@echo "🧹 Cleaning binary and coverage files..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML) $(BIN)

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

up: ## Build and start docker containers
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
