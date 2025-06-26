APP_NAME = hospital-management-system
BUILD_DIR = bin
MAIN_PATH = cmd/server/main.go
DOCKER_IMAGE = hospital-management-system:latest

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = gofmt

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/server $(MAIN_PATH)

.PHONY: run
run: build
	@echo "Starting $(APP_NAME)..."
	./$(BUILD_DIR)/server

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-watch
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -c $(GOTEST) -v ./...

.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: db-create
db-create:
	@echo "Creating database..."
	createdb hospital_management

.PHONY: db-drop
db-drop:
	@echo "Dropping database..."
	dropdb hospital_management

.PHONY: db-reset
db-reset: db-drop db-create
	@echo "Database reset complete"

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/mock_interfaces.go

.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/server-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/server-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/server-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/server-darwin-arm64 $(MAIN_PATH)

.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) -u github.com/golang/mock/mockgen@latest

.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: profile
profile:
	@echo "Running with profiling..."
	$(GOBUILD) -o $(BUILD_DIR)/server-profile $(MAIN_PATH)
	./$(BUILD_DIR)/server-profile -cpuprofile=cpu.prof -memprofile=mem.prof

.PHONY: deps-check
deps-check:
	@echo "Checking for outdated dependencies..."
	$(GOCMD) list -u -m all

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-watch     - Run tests in watch mode"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  dev-setup      - Setup development environment"
	@echo "  db-create      - Create database"
	@echo "  db-drop        - Drop database"
	@echo "  db-reset       - Reset database"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  mocks          - Generate mocks"
	@echo "  security       - Run security scan"
	@echo "  build-all      - Build for multiple platforms"
	@echo "  install-tools  - Install development tools"
	@echo "  bench          - Run benchmarks"
	@echo "  profile        - Run with profiling"
	@echo "  deps-check     - Check for outdated dependencies"
	@echo "  help           - Show this help message"

.DEFAULT_GOAL := help 