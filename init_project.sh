#!/usr/bin/env bash

set -euo pipefail

echo "🛡️  KOMGRIP PROJECT INITIALIZATION"
echo "=================================="
echo ""

run_sed() {
    local pattern="$1"
    local file="$2"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$pattern" "$file"
    else
        sed -i "$pattern" "$file"
    fi
}

read -p "Enter your project name (e.g., my-awesome-project): " PROJECT_NAME

if [ -z "$PROJECT_NAME" ]; then
    echo "❌ Project name cannot be empty!"
    exit 1
fi

PROJECT_NAME_CLEAN=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')

read -p "Enter your Go module path (e.g., github.com/yourname/$PROJECT_NAME_CLEAN): " MODULE_PATH

if [ -z "$MODULE_PATH" ]; then
    echo "❌ Module path cannot be empty!"
    exit 1
fi

echo ""
echo "📝 Configuration:"
echo "  Project Name: $PROJECT_NAME_CLEAN"
echo "  Module Path: $MODULE_PATH"
echo ""
read -p "Proceed with initialization? [Y/n]: " CONFIRM

if [[ "$CONFIRM" =~ ^[Nn]$ ]]; then
    echo "❌ Initialization aborted."
    exit 0
fi

echo ""
echo "🔄 Step 1: Renaming Go module references..."

if [ -f "api/go.mod" ]; then
    run_sed "s|module github.com/komgrip/starter-kit|module $MODULE_PATH|g" "api/go.mod"
    echo "  ✅ api/go.mod updated"
fi

if [ -d "api" ]; then
    find api -type f \( -name "*.go" -o -name "*.mod" -o -name "*.sum" \) | while IFS= read -r file; do
        run_sed "s|github.com/komgrip/starter-kit|$MODULE_PATH|g" "$file"
    done
    echo "  ✅ All Go files updated"
fi

echo ""
echo "📦 Step 2: Creating environment files..."

if [ -f ".env.example" ] && [ ! -f ".env" ]; then
    cp .env.example .env
    run_sed "s|PROJECT_NAME=komgrip|PROJECT_NAME=$PROJECT_NAME_CLEAN|g" ".env"
    echo "  ✅ .env created from .env.example"
elif [ ! -f ".env" ]; then
    cat > .env <<EOF
PROJECT_NAME=$PROJECT_NAME_CLEAN
APP_ENV=development
API_PORT=8080
WEB_PORT=3000

POSTGRES_USER=komgrip
POSTGRES_PASSWORD=komgrip_secret
POSTGRES_DB=komgrip_db
POSTGRES_PORT=5432

REDIS_PASSWORD=komgrip_secret
REDIS_PORT=6379

JWT_SECRET=$(openssl rand -base64 32)
NUXT_PUBLIC_API_BASE=http://localhost:8080
NODE_ENV=development
EOF
    echo "  ✅ .env created with defaults"
else
    echo "  ⚠️  .env already exists, skipping"
fi

if [ -f "api/.env.example" ] && [ ! -f "api/.env" ]; then
    cp api/.env.example api/.env
    echo "  ✅ api/.env created"
elif [ ! -f "api/.env" ]; then
    cat > api/.env <<EOF
APP_ENV=development
APP_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=komgrip
POSTGRES_PASSWORD=komgrip_secret
POSTGRES_DB=komgrip_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=komgrip_secret

JWT_SECRET=$(openssl rand -base64 32)
EOF
    echo "  ✅ api/.env created with defaults"
else
    echo "  ⚠️  api/.env already exists, skipping"
fi

if [ -f "web/.env.example" ] && [ ! -f "web/.env" ]; then
    cp web/.env.example web/.env
    echo "  ✅ web/.env created"
elif [ ! -f "web/.env" ]; then
    cat > web/.env <<EOF
NUXT_PUBLIC_API_BASE=http://localhost:8080
NODE_ENV=development
EOF
    echo "  ✅ web/.env created with defaults"
else
    echo "  ⚠️  web/.env already exists, skipping"
fi

echo ""
echo "🧹 Step 3: Cleaning up..."

if [ -f "api/go.sum" ]; then
    rm -f api/go.sum
    echo "  ✅ Removed api/go.sum (will be regenerated)"
fi

echo ""
echo "✅ INITIALIZATION COMPLETE!"
echo ""
echo "Next steps:"
echo "  1. Review and update .env files with your secrets"
echo "  2. Run 'make up' to start all services"
echo "  3. Run 'make migrate-up' to initialize the database"
echo "  4. Visit http://localhost:3000 for the web app"
echo "  5. Visit http://localhost:8080 for the API"
echo ""
echo "Happy coding! 🚀"
