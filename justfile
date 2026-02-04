# https://just.systems

set dotenv-load := true

# Default recipe
default: dev

# --- Main Development Commands ---

# Run the full stack in development mode (Backend + Frontend)
dev: db-start
    @echo "🚀 Starting Go Backend..."
    # We run the backend in the background.
    # Note: Ctrl+C might not kill the background process immediately in some shells.
    (cd backend && go run main.go) & \
    (cd frontend && pnpm install && pnpm dev)

# Build both backend and frontend for production
build:
    @echo "📦 Building Backend..."
    mkdir -p bin
    cd backend && go build -o ../bin/server main.go
    @echo "📦 Building Frontend..."
    cd frontend && pnpm install && pnpm build

# Clean build artifacts
clean:
    rm -rf bin
    rm -rf frontend/dist
    rm -rf frontend/node_modules

# --- Database Helpers ---

# Start the local postgres instance
db-start:
    pg-manage start || echo "⚠️ Database might already be running"

# Stop the local postgres instance
db-stop:
    pg-manage stop

# Initialize the database environment
db-init:
    pg-manage init
    pg-manage start || echo "⚠️ Database might already be running"
    pg-manage create-user-db || echo "⚠️ Database might already exist"

# Open a psql shell to the database
db-shell:
    pg-manage start || echo "⚠️ Database might already be running"
    pg-manage shell
