# CMU Insta

Automating frontend for the CMU prefrosh page.

**Tech Stack:**
- **Frontend:** Svelte (Vite + TypeScript)
- **Backend:** Go
- **Database:** PostgreSQL

### Getting Started

To get started, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmuinsta.git
   ```
2. Make sure you have NixOS with DirEnv installed, similar if not same as [Terrier contribution setup](https://github.com/ScottyLabs/terrier/blob/main/CONTRIBUTING.md)
3. Navigate to the project directory:
   ```bash
   git clone https://github.com/etashj/cmuinsta.git
   cd cmuinsta
   ```
4. Copy the `.env.sample` file and modify
5. Initialize the database
   ```bash
   just db-init
   ```
6. Run the app for development or production
   ```bash
   just dev
   just up
   ```
7. Kill the app
   ```bash
   just down
   ```

### Commands

- `just dev`: Start the full stack in development mode.
- `just build`: Compile the backend and frontend for production.
- `just clean`: Remove build artifacts (`bin/`, `dist/`, `node_modules/`).
- `just db-init`: Initialize the local Postgres database.
- `just db-start`: Start the database manually.
- `just db-stop`: Stop the database.
- `just db-shell`: Open a `psql` shell to the local database.