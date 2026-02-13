set dotenv-load := true

# Global Variables
export PGDATA           := invocation_directory() + "/.pg_data"
export PGHOST           := invocation_directory() + "/tmp"
export LOGFILE          := PGDATA + "/logfile"
export POSTS_STORE_DIR  := invocation_directory() + "/.posts_store"

# Use the correct database name here
export DB_NAME      := "cmuinsta"

# Construct the URL using the variables defined above
export DATABASE_URL := "host=" + PGHOST + " user=postgres dbname=" + DB_NAME + " sslmode=disable"

# Global Port Config (with defaults)
BACKEND_PORT  := env_var_or_default('BACKEND_PORT', '8080')
FRONTEND_PORT := env_var_or_default('FRONTEND_PORT', '5173')

# OIDC Configuration (loaded from .env)
OIDC_ISSUER_URL   := env_var_or_default('OIDC_ISSUER_URL', '')
OIDC_CLIENT_ID    := env_var_or_default('OIDC_CLIENT_ID', '')
OIDC_CLIENT_SECRET := env_var_or_default('OIDC_CLIENT_SECRET', '')
OIDC_REDIRECT_URI := env_var_or_default('OIDC_REDIRECT_URI', 'http://localhost:5173/oauth/callback')

# Admin Configuration
ADMIN_IDS := env_var_or_default('ADMIN_IDS', '')

# Default recipe
default:
    @just --list

# --- The "One Script" Solution ---

# Setup, Sync, and Run everything
up: db-start
    mkdir -p {{POSTS_STORE_DIR}}
    @echo "📦 Syncing all dependencies..."
    cd backend && go mod tidy
    cd frontend && bun install
    @echo "🚀 Starting Full Stack (Backend: {{BACKEND_PORT}}, Frontend: {{FRONTEND_PORT}})..."
    (cd backend && \
        PORT={{BACKEND_PORT}} \
        POSTS_STORE_DIR={{POSTS_STORE_DIR}} \
        OIDC_ISSUER_URL={{OIDC_ISSUER_URL}} \
        OIDC_CLIENT_ID={{OIDC_CLIENT_ID}} \
        OIDC_CLIENT_SECRET={{OIDC_CLIENT_SECRET}} \
        OIDC_REDIRECT_URI={{OIDC_REDIRECT_URI}} \
        ADMIN_IDS={{ADMIN_IDS}} \
        go run main.go) & \
    (cd frontend && bun dev --port {{FRONTEND_PORT}})

# Stop the full stack and cleanup background processes safely
down: db-stop
    @echo "🛑 Cleaning up your dev processes..."
    @# -U $(id -u) ensures you only kill processes owned by YOU
    -pkill -u $(id -u) -f "go run main.go" || true
    -pkill -u $(id -u) -f "vite" || true
    -pkill -u $(id -u) -f "bun" || true
    @echo "✅ Everything is down."

# --- Development Commands ---

dev: db-start
    mkdir -p {{POSTS_STORE_DIR}}
    @echo "📦 Syncing all dependencies..."
    cd backend && go mod tidy
    cd frontend && bun install
    @echo "🚀 Starting Full Stack (Backend: {{BACKEND_PORT}}, Frontend: {{FRONTEND_PORT}})..."
    (cd backend && \
        PORT={{BACKEND_PORT}} \
        POSTS_STORE_DIR={{POSTS_STORE_DIR}} \
        OIDC_ISSUER_URL={{OIDC_ISSUER_URL}} \
        OIDC_CLIENT_ID={{OIDC_CLIENT_ID}} \
        OIDC_CLIENT_SECRET={{OIDC_CLIENT_SECRET}} \
        OIDC_REDIRECT_URI={{OIDC_REDIRECT_URI}} \
        ADMIN_IDS={{ADMIN_IDS}} \
        go run main.go) & \
    (cd frontend && bun dev --port {{FRONTEND_PORT}})


build:
    @echo "📦 Building project..."
    mkdir -p bin
    cd backend && go build -o ../bin/server main.go
    cd frontend && bun install && bun run build

clean:
    rm -rf bin frontend/dist frontend/node_modules .pg_data tmp

# --- Database Management ---

db-start:
    @mkdir -p {{PGHOST}}
    @if [ ! -d "{{PGDATA}}" ]; then \
        echo "📦 Initializing New Database Cluster..."; \
        initdb --auth=trust --no-locale --encoding=UTF8 > /dev/null; \
        pg_ctl -D {{PGDATA}} -l {{LOGFILE}} -o "-k {{PGHOST}}" start; \
        sleep 2; \
        createuser -h {{PGHOST}} -s postgres || true; \
        createdb -h {{PGHOST}} -O postgres postgres || true; \
    fi
    @pg_ctl status -D {{PGDATA}} > /dev/null || pg_ctl -D {{PGDATA}} -l {{LOGFILE}} -o "-k {{PGHOST}}" start
    @# Verify the postgres role exists (handles cases where init was interrupted)
    @psql -h {{PGHOST}} -U $(whoami) -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='postgres'" | grep -q 1 || createuser -h {{PGHOST}} -s postgres
    @echo "🐘 Postgres is active at {{PGHOST}}"

db-stop:
    @pg_ctl -D {{PGDATA}} stop

db-reset:
    @just db-stop || true
    rm -rf {{PGDATA}} {{PGHOST}}
    @just db-start
    @just db-init

db-shell: db-start
    psql -h {{PGHOST}} -U postgres -d postgres

# Initialize the database (Run this once)
db-init:
    psql -U etashj -d postgres -c "CREATE DATABASE cmuinsta;"
