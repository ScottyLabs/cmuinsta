# CMU Insta

Automating frontend for the CMU prefrosh page.

**Tech Stack:**
- **Frontend:** Svelte (Vite + TypeScript)
- **Backend:** Go
- **Database:** PostgreSQL

### TODO
 - Add image processing on client side (compress and crops to ideal aspect ratio as supported by instagram API)
 - Image processing for the first slide, we need to be able to frame an image like the Co2030 account does in the red or blue frame randomly (it alternates iirc)
 - Post preview for client
 - Continous uploads: when a client uploads an image it should automatically sync to the backend ASAP so that their progress is saved
 - All of backend should be reviewed/rewritten

### Getting Started

To get started, follow these steps:

1. Clone the repository and move into that directory:
   ```bash
   git clone https://codeberg.org/ScottyLabs/cmuinsta
   cd cmuinsta
   ```
2. Make sure you have NixOS with DirEnv installed, similar if not same as [Terrier contribution setup](https://github.com/ScottyLabs/terrier/blob/main/CONTRIBUTING.md)
4. Copy the `.env.sample` file and modify
5. Inititalize the database (only run this once)
   ```bash
   just db-init
    ```
5. Run the app for development or production
   ```bash
   just dev
   just up
   ```
6. Kill the app
   ```bash
   just down
   ```
7. If something goes wrong idk what to do tbh

### Commands

- `just`: List all commands
- `just dev`: Start the full stack in development mode.
- `just up`: Equivalent to dev
- `just build`: Compile the backend and frontend for production.
- `just down`: Kill all processes
- `just clean`: Deletes all temp files, databases, etc. 
- `just db-init`: Initialize the local Postgres database.
- `just db-start`: Start the database manually.
- `just db-stop`: Stop the database.
- `just db-shell`: Open a `psql` shell to the local database.
- `just db-reset`: Delete `postgres` directory
