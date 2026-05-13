# CMU Insta - Developer Guide

This project uses a split-stack architecture:
- **Frontend**: Svelte (Vite + TypeScript)
- **Backend**: Go (Standard Library `net/http`)
- **Database**: PostgreSQL (Managed via Nix)

## Project Structure

- `/frontend`: Svelte application
- `/backend`: Go API server
- `/.pg_data`: Local Postgres data (managed by `pg-manage`)
- `flake.nix`: Nix development environment configuration
- `justfile`: Command runner

## Running the Project

Ensure you are in the Nix shell (`nix develop` or `direnv allow`).

```sh
# Start Database, Backend, and Frontend
just dev

# Build for production
just build
```

## Backend (Go)

The backend is a standard Go HTTP server located in `backend/`.

## Development

Refer to the RFC at ./rfcs/ to design choices at a high level and the overall 
codebase for implementation and library choices.
