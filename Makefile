.PHONY: help ensure-go-sum up down restart logs clean init shell-db shell-mongo shell-redis shell-api shell-web ps migrate-up migrate-down test-api test-web build wait-docker

PROJECT_NAME ?= komgrip

# Docker Desktop on this host returns 500 for default API (v1.51); force 1.44 so down/up work.
export DOCKER_API_VERSION ?= 1.44

ensure-go-sum:
	@if [ ! -f api/go.sum ]; then \
		echo "✨ Creating placeholder api/go.sum..."; \
		touch api/go.sum; \
		echo "✅ api/go.sum created."; \
	fi

help:
	@echo "🛡️  KOMGRIP GOD-TIER MAKEFILE"
	@echo ""
	@echo "Available commands:"
	@echo "  make up              - Start all services (detached)"
	@echo "  make run             - Start all services in terminal (see logs, Ctrl+C to stop)"
	@echo "  make down            - Stop all services"
	@echo "  make restart         - Restart all services"
	@echo "  make logs            - Show logs (all services)"
	@echo "  make logs-api        - Show API service logs only"
	@echo "  make clean           - Stop services and remove volumes (DESTRUCTIVE)"
	@echo "  make init            - Initialize project (rename, setup .env)"
	@echo "  make ps              - List running containers"
	@echo ""
	@echo "Shell access:"
	@echo "  make shell-db        - PostgreSQL shell (psql)"
	@echo "  make shell-mongo     - MongoDB shell (mongosh)"
	@echo "  make shell-redis     - Redis CLI"
	@echo "  make shell-api       - API container shell"
	@echo "  make shell-web       - Web container shell"
	@echo ""
	@echo "Development:"
	@echo "  make build           - Build Docker images"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback migrations"
	@echo "  make test-api        - Run API tests"
	@echo "  make test-web        - Run Web tests"

# Wait for Docker daemon (use before build/up if Docker Desktop just started)
wait-docker:
	@echo "⏳ Waiting for Docker daemon..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do \
		if docker info >/dev/null 2>&1; then echo "✅ Docker is ready."; exit 0; fi; \
		echo "   attempt $$i/20..."; sleep 6; \
	done; \
	echo "❌ Docker did not respond. Start Docker Desktop and run make again."; exit 1

build: ensure-go-sum wait-docker
	@echo "🔨 Building Docker images..."
	docker-compose build
	@echo "✅ Build complete."

up: ensure-go-sum wait-docker
	@echo "🚀 Starting $(PROJECT_NAME) services..."
	docker-compose up -d
	@echo "✅ Services started. Use 'make logs' to view logs."

# Run all services in foreground (logs in terminal). Ctrl+C to stop.
run: ensure-go-sum wait-docker
	@echo "🚀 Starting all services (logs in terminal, Ctrl+C to stop)..."
	docker-compose up

down:
	@echo "🛑 Stopping $(PROJECT_NAME) services..."
	@docker-compose down || (echo "⚠️  Docker not available or already stopped."; exit 0)
	@echo "✅ Services stopped."

restart: wait-docker down up

logs:
	@echo "📋 Showing logs for all services (Ctrl+C to exit)..."
	docker-compose logs -f

logs-api:
	@echo "📋 API service logs (Ctrl+C to exit)..."
	docker-compose logs -f api

clean:
	@echo "⚠️  WARNING: This will DELETE all volumes and data!"
	@read -p "Are you sure? [y/N]: " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "🗑️  Cleaning up..."; \
		docker-compose down -v --remove-orphans; \
		echo "✅ Cleanup complete."; \
	else \
		echo "❌ Aborted."; \
	fi

init:
	@echo "🔧 Initializing project..."
	@if [ -f ./init_project.sh ]; then \
		chmod +x ./init_project.sh; \
		./init_project.sh; \
	else \
		echo "❌ init_project.sh not found!"; \
		exit 1; \
	fi

ps:
	@echo "📊 Running containers:"
	docker-compose ps

shell-db:
	@echo "🐘 Connecting to PostgreSQL..."
	docker-compose exec postgres psql -U $${POSTGRES_USER:-komgrip} -d $${POSTGRES_DB:-komgrip_db}

shell-mongo:
	@echo "🍃 Connecting to MongoDB..."
	docker-compose exec mongo mongosh -u $${MONGO_USER:-komgrip} -p $${MONGO_PASSWORD:-komgrip_secret} --authenticationDatabase admin

shell-redis:
	@echo "⚡ Connecting to Redis..."
	docker-compose exec redis redis-cli -a $${REDIS_PASSWORD:-komgrip_secret}

shell-api:
	@echo "🔧 Accessing API container shell..."
	docker-compose exec api /bin/sh

shell-web:
	@echo "🌐 Accessing Web container shell..."
	docker-compose exec web /bin/sh

migrate-up:
	@echo "📈 Running database migrations..."
	docker-compose exec api go run cmd/migrate/main.go up

migrate-down:
	@echo "📉 Rolling back migrations..."
	docker-compose exec api go run cmd/migrate/main.go down

test-api:
	@echo "🧪 Running API tests..."
	docker-compose exec api go test -v -race -coverprofile=coverage.out ./...
	@echo "📊 Coverage report:"
	docker-compose exec api go tool cover -func=coverage.out

test-web:
	@echo "🧪 Running Web tests..."
	docker-compose exec web npm run test
