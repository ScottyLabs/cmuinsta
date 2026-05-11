set dotenv-load := true

# Global Variables
export PGDATA          := invocation_directory() + "/.pg_data"
export PGHOST          := invocation_directory() + "/.pg_tmp"
export LOGFILE         := PGDATA + "/logfile"
export POSTS_STORE_DIR := invocation_directory() + "/.posts_store"

# Database
export DB_NAME     := "cmuinsta"
export DATABASE_URL := "postgresql://postgres@/" + DB_NAME + "?host=" + PGHOST

# Ports (overridable via .env)
BACKEND_PORT  := env_var_or_default('BACKEND_PORT', '8080')
FRONTEND_PORT := env_var_or_default('FRONTEND_PORT', '5173')

# OIDC Configuration (loaded from .env)
OIDC_ISSUER_URL    := env_var_or_default('OIDC_ISSUER_URL', '')
OIDC_CLIENT_ID     := env_var_or_default('OIDC_CLIENT_ID', '')
OIDC_CLIENT_SECRET := env_var_or_default('OIDC_CLIENT_SECRET', '')
OIDC_REDIRECT_URI  := env_var_or_default('OIDC_REDIRECT_URI', 'http://localhost:' + FRONTEND_PORT + '/oauth/callback')

# Admin Configuration
ADMIN_IDS := env_var_or_default('ADMIN_IDS', '')

# Default recipe
default:
    @just --list

# --- Main Commands ---

# Install dependencies, start DB, and run the full stack with hot reload
dev: db-ensure
    @mkdir -p {{POSTS_STORE_DIR}}
    @echo "📦 Syncing dependencies..."
    @cd backend && go mod tidy
    @cd frontend && bun install
    @echo "🚀 Starting full stack (backend: {{BACKEND_PORT}}, frontend: {{FRONTEND_PORT}})..."
    @trap 'kill $(jobs -p) 2>/dev/null; just db-stop' EXIT; \
        (cd backend && \
            DATABASE_URL={{DATABASE_URL}} \
            PORT={{BACKEND_PORT}} \
            POSTS_STORE_DIR={{POSTS_STORE_DIR}} \
            OIDC_ISSUER_URL={{OIDC_ISSUER_URL}} \
            OIDC_CLIENT_ID={{OIDC_CLIENT_ID}} \
            OIDC_CLIENT_SECRET={{OIDC_CLIENT_SECRET}} \
            OIDC_REDIRECT_URI={{OIDC_REDIRECT_URI}} \
            ADMIN_IDS={{ADMIN_IDS}} \
            go run main.go) & \
        (cd frontend && bun dev --port {{FRONTEND_PORT}})

# Build production binaries
build: db-ensure
    @echo "📦 Building..."
    @mkdir -p bin
    @cd backend && go build -o ../bin/server main.go
    @cd frontend && bun install && bun run build
    @echo "✅ Build complete → bin/server and frontend/dist"

# Kill any running dev processes
down: db-stop
    @echo "🛑 Stopping dev processes..."
    @pkill -u $$(id -u) -f "go run main.go" || true
    @pkill -u $$(id -u) -f "vite"           || true
    @pkill -u $$(id -u) -f "bun"            || true
    @echo "✅ Done."

# Remove all build artifacts and local state (does not touch .env)
clean: down
    @echo "🧹 Cleaning..."
    @rm -rf bin frontend/dist frontend/node_modules .pg_data .pg_tmp .posts_store
    @echo "✅ Clean."

# --- Database ---

# Start Postgres if not already running, initialize cluster and DB if needed
db-ensure:
    @mkdir -p {{PGHOST}}
    @if [ ! -d "{{PGDATA}}" ]; then \
        echo "📦 Initializing Postgres cluster..."; \
        initdb --auth=trust --no-locale --encoding=UTF8 -D {{PGDATA}} > /dev/null; \
        pg_ctl -D {{PGDATA}} -l {{LOGFILE}} -o "-k {{PGHOST}}" start; \
        sleep 2; \
        createuser -h {{PGHOST}} -s postgres || true; \
        createdb   -h {{PGHOST}} -U postgres postgres || true; \
        createdb   -h {{PGHOST}} -U postgres {{DB_NAME}} || true; \
        echo "✅ Cluster ready."; \
    else \
        pg_ctl status -D {{PGDATA}} > /dev/null || \
            pg_ctl -D {{PGDATA}} -l {{LOGFILE}} -o "-k {{PGHOST}}" start; \
    fi
    @echo "🐘 Postgres is ready at {{PGHOST}}"

# Stop the Postgres server
db-stop:
    @pg_ctl -D {{PGDATA}} stop -m fast 2>/dev/null || true

# Open a psql shell
db-shell: db-ensure
    psql -h {{PGHOST}} -U postgres -d {{DB_NAME}}

# Wipe and reinitialize the database cluster from scratch
db-reset: down
    @echo "⚠️  Resetting database..."
    @pg_ctl -D {{PGDATA}} stop -m fast 2>/dev/null || true
    @rm -rf {{PGDATA}} {{PGHOST}}
    @just db-ensure
    @echo "✅ Database reset."
