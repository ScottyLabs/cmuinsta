{
  description = "A Nix-flake-based dev environment for Postgres, Go, and Svelte (via Bun)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # --- Go ---
            go
            gopls

            # --- Frontend (Bun) ---
            bun
            # We keep nodejs as a fallback, as some Svelte tools expect it
            nodejs_20

            # --- Database ---
            postgresql_16

            # --- Utilities ---
            jq
            just
          ];

          shellHook = ''
            echo "🚀 Dev Shell Active: Bun + Go + Postgres"

            # --- Local Postgres Config ---
            export PGDATA="$PWD/.pg_data"
            export PGHOST="$PWD/tmp"
            mkdir -p $PGHOST

            if [ ! -d "$PGDATA" ]; then
              echo "📦 Initializing local Postgres data..."
              initdb --auth=trust --no-locale --encoding=UTF8 > /dev/null
            fi

            echo "----------------------------------------"
            echo "🔥 BUN: use 'bun install' and 'bun dev'"
            echo "🐹 GO:  use 'go run main.go'"
            echo "🐘 DB:  'pg_ctl -l $PGDATA/logfile start'"
            echo "----------------------------------------"
          '';
        };
      }
    );
}
