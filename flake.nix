{
  description = "Shellify - Manage development workspace sessions with terminal multiplexers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { nixpkgs, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              # Go toolchain
              go
              golangci-lint
              goreleaser

              # Node for GUI
              nodejs

              # Terminal multiplexers for testing
              tmux
              zellij

              # Build tools
              gnumake
            ];

            shellHook = ''
              echo "Shellify development environment"
              echo "Go: $(go version)"
              echo "Node: $(node --version)"
              echo ""
              echo "Commands:"
              echo "  make build        - Build CLI"
              echo "  make test         - Run tests"
              echo "  make lint         - Run linter"
              echo "  make gui-dev      - Start GUI dev server"
            '';
          };
        });
    };
}
