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

### Dependencies
Dependencies are managed via `go.mod`.
- `github.com/lib/pq`: Postgres driver

### Routing
We use the standard library `http.ServeMux` (Go 1.22+).

```go
func main() {
    mux := http.NewServeMux()
    
    // Simple handler
    mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
        // ...
    })
    
    // ...
}
```

### Database
The database connection string is provided via `DATABASE_URL` environment variable.

```go
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatal(err)
}
```

## Frontend (Svelte)

The frontend is a Vite-powered Svelte application located in `frontend/`.

### Components
Components use `.svelte` files.

```html
<script lang="ts">
  export let name: string = "World"; // Prop
  let count = 0; // State

  function increment() {
    count += 1;
  }
</script>

<h1>Hello {name}!</h1>
<button on:click={increment}>Count is {count}</button>

<style>
  h1 { color: #ff3e00; }
</style>
```

### API Calls
The Vite dev server is configured (in `vite.config.js`) to proxy `/api` requests to the Go backend (defaulting to port 8080).

```ts
// In Svelte components
onMount(async () => {
    const res = await fetch("/api/data");
    const data = await res.json();
    console.log(data);
});
```

## Database Management

The Postgres instance is local to this project folder, storing data in `.pg/`.

- **Init**: `pg-manage init` (Initializes data directory)
- **Start**: `pg-manage start` (Starts the postgres process)
- **Stop**: `pg-manage stop` (Stops the postgres process)
- **Shell**: `pg-manage shell` (Opens `psql` connected to the local db)

The `just dev` command automatically handles starting the database.
