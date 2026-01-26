# https://just.systems

set dotenv-load

# Serve the application with hot reloading
dev: db-start
    dx serve

# Build the application for release
build:
    dx build --release

# Host the application (run release build)
up: db-start
    cargo run --release

# Kill the application
down: db-stop
      cargo kill

# --- Database Helpers ---

# Start the local postgres instance
db-start:
    # We ignore errors here in case it's already running
    pg-manage start || echo "Database might already be running"

# Stop the local postgres instance
db-stop:
    pg-manage stop

# Initialize the database (run once)
db-init:
    pg-manage init
    pg-manage create-user-db
