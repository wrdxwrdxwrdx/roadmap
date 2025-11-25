.PHONY: up up-dev down build rebuild logs ps clean db-shell db-tables db-describe db-size db-tables-size db-info restart-api logs-api logs-db logs-frontend wait-health test test-short test-verbose test-coverage test-unit test-integration lint lint-fix format workflow frontend-build-docker frontend-dev-docker frontend-restart frontend-logs frontend-shell frontend-clean-docker help

COVERAGE_THRESHOLD ?= 50.0

up:
	docker-compose up -d --build
	@echo "Waiting for services to be healthy..."
	@make wait-health
	@echo "✓ All services are up!"
	@bash -c 'FRONTEND_PORT=$${FRONTEND_PORT:-3000}; API_PORT=$${API_PORT:-8080}; echo "  - Frontend: http://localhost:$$FRONTEND_PORT"; echo "  - API: http://localhost:$$API_PORT"'

# Start all services with frontend in dev mode (with hot reload)
up-dev:
	docker-compose up -d postgres api
	@echo "Waiting for API to be healthy..."
	@make wait-health
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d frontend
	@echo "✓ All services are up in development mode!"
	@bash -c 'FRONTEND_PORT=$${FRONTEND_PORT:-3000}; API_PORT=$${API_PORT:-8080}; echo "  - Frontend (dev): http://localhost:$$FRONTEND_PORT"; echo "  - API: http://localhost:$$API_PORT"'

# Start all services
up-build:
	docker-compose up -d --build
	@make up

# Wait for API health check to pass
wait-health:
	@timeout=120; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose ps api | grep -q "healthy"; then \
			echo "✓ API is healthy and ready!"; \
			docker-compose ps; \
			exit 0; \
		fi; \
		echo "Waiting for API health check... ($$timeout seconds remaining)"; \
		sleep 10; \
		timeout=$$((timeout-10)); \
	done; \
	echo "⚠ Health check timeout. Check logs with 'make logs-api'"; \
	docker-compose ps; \
	exit 1

# Stop all services
down:
	docker-compose down
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down 2>/dev/null || true

# Build images
build:
	docker-compose build
	@echo "To build with dev frontend, run: docker-compose -f docker-compose.yml -f docker-compose.dev.yml build"

# Rebuild images without cache
rebuild:
	docker-compose build --no-cache

# View logs
logs:
	docker-compose logs -f

# View container status
ps:
	docker-compose ps

# Stop and remove containers and volumes
clean:
	docker-compose down -v

# Connect to database
db-shell:
	docker-compose exec postgres psql -U postgres -d roadmap

# Show all tables
db-tables:
	@docker-compose exec postgres psql -U postgres -d roadmap -c "\dt"

# Describe table structure (usage: make db-describe TABLE=users)
db-describe:
	@if [ -z "$(TABLE)" ]; then \
		echo "Usage: make db-describe TABLE=table_name"; \
		echo "Example: make db-describe TABLE=users"; \
	else \
		docker-compose exec postgres psql -U postgres -d roadmap -c "\d $(TABLE)"; \
	fi

# Show database sizes
db-size:
	@docker-compose exec postgres psql -U postgres -d roadmap -c "SELECT pg_size_pretty(pg_database_size('roadmap')) AS database_size;"

# Show table sizes
db-tables-size:
	@docker-compose exec postgres psql -U postgres -d roadmap -c "SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"

# Show database info
db-info:
	@echo "=== Database Information ==="
	@docker-compose exec postgres psql -U postgres -d roadmap -c "SELECT version();"
	@echo ""
	@echo "=== Current Database ==="
	@docker-compose exec postgres psql -U postgres -d roadmap -c "SELECT current_database();"
	@echo ""
	@echo "=== All Tables ==="
	@docker-compose exec postgres psql -U postgres -d roadmap -c "\dt"

# Restart API service
restart-api:
	docker-compose restart api

# View API logs only
logs-api:
	docker-compose logs -f api

# View PostgreSQL logs only
logs-db:
	docker-compose logs -f postgres

# View Frontend logs only
logs-frontend:
	docker-compose logs -f frontend || docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f frontend

# Frontend Docker commands
# Build frontend Docker image for production
frontend-build-docker:
	@echo "Building frontend Docker image for production..."
	docker-compose build frontend
	@echo "✓ Frontend Docker image built successfully!"
	@bash -c 'FRONTEND_PORT=$${FRONTEND_PORT:-3000}; echo "  - Frontend will be available at: http://localhost:$$FRONTEND_PORT"; echo "  - Run '\''make up'\'' to start all services"'

# Start frontend in development mode (with hot reload)
frontend-dev-docker:
	@echo "Starting frontend in development mode with Docker..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d frontend
	@bash -c 'FRONTEND_PORT=$${FRONTEND_PORT:-3000}; echo "Frontend dev server will be available at http://localhost:$$FRONTEND_PORT"'

# Restart frontend service
frontend-restart:
	@echo "Restarting frontend service..."
	docker-compose restart frontend || docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart frontend

# View frontend logs
frontend-logs:
	@make logs-frontend

# Execute command in frontend container
frontend-shell:
	docker-compose exec frontend sh || docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec frontend sh

# Clean frontend Docker images and containers
frontend-clean-docker:
	@echo "Cleaning frontend Docker artifacts..."
	docker-compose down frontend 2>/dev/null || true
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down frontend 2>/dev/null || true
	docker rmi roadmap-frontend roadmap-frontend-dev 2>/dev/null || true

# Run all tests
test:
	@echo "Running all tests..."
	@cd backend && go test ./... -v || (echo "Some tests failed. Check output above for details." && exit 1)

# Run tests and show short summary (one line)
test-short:
	@cd backend && \
	go test ./... -json -coverprofile=coverage.out 2>&1 | \
	awk 'BEGIN {passed=0; failed=0} \
		/"Action":"pass"/ {passed++} \
		/"Action":"fail"/ {failed++} \
		END { \
			total=passed+failed; \
			if (failed > 0) { \
				printf "Tests: %d/%d passed, %d failed", passed, total, failed; \
				exit 1; \
			} else if (total > 0) { \
				printf "Tests: %d/%d passed", passed, total; \
			} else { \
				printf "No tests found"; \
			} \
		}' && \
	coverage=$$(go tool cover -func=coverage.out 2>/dev/null | tail -1 | awk '{print $$3}' | sed 's/%//') && \
	if [ -n "$$coverage" ]; then \
		coverage_num=$$(echo "$$coverage" | awk '{print $$1}'); \
		threshold=$(COVERAGE_THRESHOLD); \
		if [ -n "$$coverage_num" ] && [ -n "$$threshold" ]; then \
			status=$$(echo "$$coverage_num >= $$threshold" | bc -l 2>/dev/null || echo "0"); \
			if [ "$$status" = "1" ]; then \
				echo " | Coverage: $$coverage% (target: $$threshold% ✓)"; \
			else \
				echo " | Coverage: $$coverage% (target: $$threshold% ✗)"; \
			fi; \
		else \
			echo " | Coverage: $$coverage% (target: $$threshold%)"; \
		fi; \
	else \
		echo ""; \
	fi

# Run tests with verbose output
test-verbose:
	@echo "Running all tests with verbose output..."
	@cd backend && go test ./... -v -count=1

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@cd backend && go test ./... -coverprofile=coverage.out
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

# Run only unit tests (excluding integration tests)
test-unit:
	@echo "Running unit tests..."
	@cd backend && go test ./... -v -short

# Run integration tests (requires database)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: Integration tests require TEST_DB_DSN environment variable"
	@cd backend && go test ./... -v -run Integration

# Run linters
lint:
	@echo "Running linters..."
	@cd backend && \
	if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run --timeout=5m || exit 1; \
		echo "✓ golangci-lint passed"; \
	else \
		echo "⚠ golangci-lint not found. Skipping... (Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"; \
	fi
	@cd backend && echo "Running go vet..." && go vet ./... || exit 1
	@cd backend && echo "Checking formatting..." && \
	if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "❌ Code is not formatted. Run 'make format' to fix."; \
		gofmt -s -d .; \
		exit 1; \
	fi
	@echo "✓ All linters passed"

# Fix linting issues automatically
lint-fix:
	@echo "Fixing linting issues..."
	@cd backend && \
	if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --timeout=5m; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	@cd backend && gofmt -s -w .
	@echo "✓ Linting issues fixed"

# Format code
format:
	@echo "Formatting code..."
	@cd backend && gofmt -s -w .
	@cd backend && goimports -w . || echo "goimports not found, skipping import formatting"
	@echo "✓ Code formatted"

# Run full CI workflow locally
workflow:
	@echo "========================================="
	@echo "Running CI workflow locally..."
	@echo "========================================="
	@echo ""
	@echo "Step 1/4: Running linters..."
	@$(MAKE) lint || (echo "❌ Linting failed" && exit 1)
	@echo ""
	@echo "Step 2/4: Running tests..."
	@$(MAKE) test-short || (echo "❌ Tests failed" && exit 1)
	@echo ""
	@echo "Step 3/4: Building..."
	@cd backend && go build -v -o api ./cmd/api || (echo "❌ Build failed" && exit 1)
	@echo "✓ Build successful"
	@echo ""
	@echo "Step 4/4: Checking binary..."
	@if [ -f ./backend/api ]; then \
		echo "✓ Binary exists"; \
		ls -lh ./backend/api; \
		rm -f ./backend/api; \
		echo "✓ Binary cleaned up"; \
	else \
		echo "❌ Binary not found"; \
		exit 1; \
	fi
	@echo ""
	@echo "========================================="
	@echo "✓ CI workflow completed successfully!"
	@echo "========================================="
	@echo ""
	@echo "Note: To install golangci-lint for full linting:"
	@echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
# Show help message with all available commands
help:
	@echo "Available commands:"
	@echo ""
	@echo "Backend/Docker commands:"
	@echo "  make up              - Start all Docker services (production)"
	@echo "  make up-dev          - Start all services with frontend in dev mode (hot reload)"
	@echo "  make down            - Stop all Docker services"
	@echo "  make build           - Build Docker images"
	@echo "  make rebuild         - Rebuild Docker images without cache"
	@echo "  make logs            - View all logs"
	@echo "  make logs-api        - View API logs only"
	@echo "  make logs-db         - View PostgreSQL logs only"
	@echo "  make logs-frontend   - View Frontend logs only"
	@echo "  make ps              - View container status"
	@echo "  make clean           - Stop and remove containers and volumes"
	@echo "  make restart-api     - Restart API service"
	@echo ""
	@echo "Database commands:"
	@echo "  make db-shell        - Connect to PostgreSQL database"
	@echo "  make db-tables       - Show all tables"
	@echo "  make db-describe     - Describe table (usage: make db-describe TABLE=users)"
	@echo "  make db-size         - Show database size"
	@echo "  make db-tables-size  - Show table sizes"
	@echo "  make db-info         - Show database information"
	@echo ""
	@echo "Testing commands:"
	@echo "  make test            - Run all tests"
	@echo "  make test-short      - Run tests with short summary"
	@echo "  make test-verbose    - Run tests with verbose output"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make test-unit       - Run unit tests only"
	@echo "  make test-integration - Run integration tests"
	@echo ""
	@echo "Code quality commands:"
	@echo "  make lint            - Run linters"
	@echo "  make lint-fix        - Fix linting issues automatically"
	@echo "  make format          - Format code"
	@echo "  make workflow        - Run full CI workflow locally"
	@echo ""
	@echo "Frontend Docker commands:"
	@echo "  make frontend-build-docker  - Build frontend Docker image for production"
	@echo "  make frontend-dev-docker    - Start frontend in development mode (Docker, hot reload)"
	@echo "  make frontend-restart       - Restart frontend service"
	@echo "  make frontend-logs          - View frontend logs"
	@echo "  make frontend-shell         - Open shell in frontend container"
	@echo "  make frontend-clean-docker  - Clean frontend Docker images and containers"
	@echo ""
	@echo "  make help            - Show this help message"
