# https://just.systems

set dotenv-load:=true

DB_NAME     := env_var('DB_NAME')
ADMIN_EMAILS := env_var('ADMIN_EMAILS')
ADMIN_GROUP := "ADMINS"

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

# Initialize DB with Group-based authentication
db-init:
    pg-manage init
    pg-manage start || echo "Database might already be running"
    @echo "Creating database: {{DB_NAME}}..."
    psql -d postgres -c "CREATE DATABASE {{DB_NAME}};" || echo "Database already exists."

    @echo "Setting up Group Role: {{ADMIN_GROUP}}..."
    # Create the group role (NOLOGIN means you can't log in as the group itself)
    psql -d {{DB_NAME}} -c "CREATE ROLE {{ADMIN_GROUP}} WITH NOLOGIN;" || echo "Group already exists."

    # Grant permissions to the Group
    psql -d {{DB_NAME}} -c "GRANT ALL PRIVILEGES ON DATABASE {{DB_NAME}} TO {{ADMIN_GROUP}};"
    psql -d {{DB_NAME}} -c "GRANT ALL PRIVILEGES ON SCHEMA public TO {{ADMIN_GROUP}};"
    # Ensure future tables created by one admin are accessible by others
    psql -d {{DB_NAME}} -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO {{ADMIN_GROUP}};"

    @for user in {{ADMIN_EMAILS}}; do \
        echo "Adding $$user to {{ADMIN_GROUP}}..."; \
        psql -d postgres -c "CREATE USER \"$$user\";" || echo "User $$user already exists."; \
        psql -d {{DB_NAME}} -c "GRANT {{ADMIN_GROUP}} TO \"$$user\";"; \
    done
    @echo "Group-based initialization complete."

db-shell:
    pg-manage start || echo "Database might already be running"
    pg-manage shell
