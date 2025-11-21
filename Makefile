.PHONY: up down build rebuild logs ps clean db-shell db-tables db-describe db-size db-tables-size db-info restart-api logs-api logs-db wait-health test test-short test-verbose test-coverage test-unit test-integration

# Coverage threshold (minimum required coverage percentage)
COVERAGE_THRESHOLD ?= 65.0

# Start all services
up:
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@make wait-health

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

# Build images
build:
	docker-compose build

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

