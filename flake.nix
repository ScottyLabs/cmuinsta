{
  description = "Rust Dioxus Full Stack with Postgres";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, rust-overlay, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        # Rust toolchain setup
        # Includes wasm32-unknown-unknown for Dioxus web client compilation
        rustToolchain = pkgs.rust-bin.stable.latest.default.override {
          extensions = [ "rust-src" "rust-analyzer" ];
          targets = [ "wasm32-unknown-unknown" ];
        };

        # Helper script to manage a local Postgres instance in .pg/
        pgScript = pkgs.writeShellScriptBin "pg-manage" ''
          export PGDATA=$PWD/.pg/data
          export PGHOST=$PWD/.pg/socket
          export LOGFILE=$PWD/.pg/logfile

          mkdir -p $PGHOST

          case "$1" in
            init)
              if [ ! -d $PGDATA ]; then
                echo "Initializing Postgres..."
                ${pkgs.postgresql}/bin/initdb -D $PGDATA --no-locale --encoding=UTF8

                # Configure to use local socket only
                echo "unix_socket_directories = '$PGHOST'" >> $PGDATA/postgresql.conf
                echo "listen_addresses = '''" >> $PGDATA/postgresql.conf

                echo "Postgres initialized in .pg/"
              else
                echo "Postgres already initialized."
              fi
              ;;
            start)
              if [ ! -d $PGDATA ]; then
                echo "Run 'pg-manage init' first."
                exit 1
              fi
              ${pkgs.postgresql}/bin/pg_ctl -D $PGDATA -l $LOGFILE -o "-k $PGHOST" start
              echo "Postgres started. Connect via socket at $PGHOST"
              ;;
            stop)
               ${pkgs.postgresql}/bin/pg_ctl -D $PGDATA stop
               ;;
            create-user-db)
               # creates a db with the same name as the current user
               ${pkgs.postgresql}/bin/createdb -h $PGHOST $USER
               ;;
            shell)
               ${pkgs.postgresql}/bin/psql -h $PGHOST postgres
               ;;
            *)
              echo "Usage: pg-manage {init|start|stop|create-user-db|shell}"
              ;;
          esac
        '';

      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Rust environment
            rustToolchain
            trunk        # Required for building Dioxus web
            cargo-watch

            # Database
            postgresql
            pgScript     # Our custom management script

            # System libs
            pkg-config
            openssl
          ];

          # Environment setup
          shellHook = ''
            export PGDATA=$PWD/.pg/data
            export PGHOST=$PWD/.pg/socket

            # Common Env var for SQLx and other tools
            # Note: The host part points to the unix socket directory
            export DATABASE_URL="postgresql:///$USER?host=$PGHOST"

            echo "🚀 Rust Dioxus + Postgres Dev Environment"
            echo "Commands available:"
            echo "  pg-manage init           -> Initialize database in .pg/"
            echo "  pg-manage start          -> Start the database server"
            echo "  pg-manage create-user-db -> Create a DB named '$USER'"
            echo "  trunk serve              -> Run your Dioxus app"
          '';
        };
      }
    );
}
