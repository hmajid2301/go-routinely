{
  description = "Development environment for Goroutinely";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      ...
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };

        myPackages = with pkgs; [
          go

          goose
          air
          golangci-lint
          gotools
          gotestsum
          gocover-cobertura
          go-task
          go-mockery

          tailwindcss
          sqlc
          concurrently

          rustywind
        ];

        devShellPackages =
          with pkgs;
          myPackages
          ++ [
            gomod2nix.packages.${system}.default
          ];
      in
      rec {
        packages.default = pkgs.buildGoApplication {
          pname = "goroutinely";
          version = "0.1.0";
          src = ./.;
          modules = ./gomod2nix.toml;

          # Build Tailwind CSS before building Go binary
          preBuild = ''
            ${pkgs.tailwindcss}/bin/tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css --minify
          '';
        };

        packages.container = pkgs.dockerTools.buildImage {
          name = "goroutinely";
          tag = "latest";
          created = "now";
          copyToRoot = pkgs.buildEnv {
            name = "image-root";
            paths = [
              packages.default
              pkgs.cacert
            ];
            pathsToLink = [ "/bin" ];
          };
          config = {
            ExposedPorts = {
              "8080/tcp" = { };
            };
            Cmd = [ "${packages.default}/bin/goroutinely" ];
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              "SSL_CERT_DIR=${pkgs.cacert}/etc/ssl/certs/"
            ];
            Labels = {
              service = "goroutinely";
            };
          };
        };

        devShells.default = pkgs.mkShell {
          packages = devShellPackages;
          shellHook = ''
            export GOOSE_DRIVER=postgres
            export GOOSE_DBSTRING="''${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/goroutinely?sslmode=disable}"
            export GOOSE_MIGRATION_DIR="internal/store/db/migrations"

            echo "🌱 Goroutinely development environment"
            echo "Available commands:"
            echo "  goose up         - Run database migrations"
            echo "  goose status     - Check migration status"
            echo "  air              - Start development server with live reload"
            echo "  sqlc generate    - Generate type-safe SQL code"
            echo "  task dev         - Start full dev environment"
          '';
        };
      }
    ));
}
