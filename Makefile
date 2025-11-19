.PHONY: up down build rebuild logs ps clean db-shell restart-api logs-api logs-db wait-health

# Start all services
up:
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@make wait-health

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
		sleep 5; \
		timeout=$$((timeout-5)); \
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

# Restart API service
restart-api:
	docker-compose restart api

# View API logs only
logs-api:
	docker-compose logs -f api

# View PostgreSQL logs only
logs-db:
	docker-compose logs -f postgres

