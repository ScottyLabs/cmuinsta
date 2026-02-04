# CMU Insta

Automating frontend for the CMU prefrosh page.

**Tech Stack:**
- **Frontend:** Svelte (Vite + TypeScript)
- **Backend:** Go
- **Database:** PostgreSQL

### Getting Started

To get started, follow these steps:

<<<<<<< HEAD
1. **Clone the repository:**
=======
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmuinsta.git
   ```
2. Make sure you have NixOS with DirEnv installed, similar if not same as [Terrier contribution setup](https://github.com/ScottyLabs/terrier/blob/main/CONTRIBUTING.md)
3. Navigate to the project directory:
>>>>>>> refs/remotes/origin/main
   ```bash
   git clone https://github.com/etashj/cmuinsta.git
   cd cmuinsta
   ```
<<<<<<< HEAD

2. **Environment Setup:**
   Ensure you have [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/) installed.
   ```bash
   direnv allow
   ```
   This will automatically install Go, Node.js, Postgres, and other tools.

3. **Initialize the Database:**
   ```bash
   just db-init
   ```

4. **Run Development Server:**
   This starts the Postgres database, the Go backend, and the Svelte frontend concurrently.
   ```bash
   just dev
=======
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
>>>>>>> refs/remotes/origin/main
   ```

### Commands

- `just dev`: Start the full stack in development mode.
- `just build`: Compile the backend and frontend for production.
- `just clean`: Remove build artifacts (`bin/`, `dist/`, `node_modules/`).
- `just db-init`: Initialize the local Postgres database.
- `just db-start`: Start the database manually.
- `just db-stop`: Stop the database.
- `just db-shell`: Open a `psql` shell to the local database.