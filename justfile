set dotenv-load := true

# Global Variables
export PGDATA   := invocation_directory() + "/.pg_data"
export PGHOST   := invocation_directory() + "/tmp"
export LOGFILE  := PGDATA + "/logfile"
# Force absolute path for DATABASE_URL so 'cd backend' doesn't break the socket path
export DATABASE_URL := "host=" + PGHOST + " user=postgres dbname=postgres sslmode=disable"

# Default recipe
default:
    @just --list

# --- The "One Script" Solution ---

# Setup, Sync, and Run everything
up: db-start
    @echo "📦 Syncing all dependencies..."
    cd backend && go mod tidy
    cd frontend && bun install
    @echo "🚀 Starting Full Stack..."
    (cd backend && go run main.go) & (cd frontend && bun dev)

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
    @echo "🚀 Starting Full Stack..."
    (cd backend && go run main.go) & (cd frontend && bun dev)

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

db-shell: db-start
    psql -h {{PGHOST}} -U postgres -d postgres
