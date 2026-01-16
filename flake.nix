{
  description = "Shellify - Manage development workspace sessions with terminal multiplexers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs systems;

      # Version from file
      version = builtins.replaceStrings [ "\n" ] [ "" ] (builtins.readFile ./VERSION);

      # Package builder function
      mkShellify = pkgs:
        let
          # Build GUI first with buildNpmPackage
          gui = pkgs.buildNpmPackage {
            pname = "shellify-gui";
            inherit version;

            src = ./gui;

            npmDepsHash = "sha256-1c5l1MZLXLIZuAi/xkPIf3etly5OF2khXcvynTfIJMM=";

            nodejs = pkgs.nodejs_22;

            buildPhase = ''
              runHook preBuild
              npm run build
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p $out
              cp -r dist/* $out/
              runHook postInstall
            '';
          };
        in
        pkgs.buildGoModule {
          pname = "shellify";
          inherit version;

          src = ./.;

          vendorHash = "sha256-UFLy0O2HuxlnX3WcrmOa1MR5dh2hlaZnyqFDI/GJajA=";

          nativeBuildInputs = with pkgs; [
            installShellFiles
          ];

          # Copy pre-built GUI to gui/dist before Go build
          preBuild = ''
            mkdir -p gui/dist
            cp -r ${gui}/* gui/dist/
          '';

          # Build with embedded GUI
          tags = [ "embed_gui" ];

          ldflags = [
            "-s"
            "-w"
            "-X github.com/ajmasia/shellify/internal/interfaces/cli.Version=${version}"
          ];

          doCheck = false;

          postInstall = ''
            installShellCompletion --cmd sfy \
              --bash <($out/bin/sfy completion bash) \
              --zsh <($out/bin/sfy completion zsh) \
              --fish <($out/bin/sfy completion fish)
          '';

          meta = with pkgs.lib; {
            description = "Terminal multiplexer session manager for tmux and zellij";
            homepage = "https://github.com/ajmasia/shellify";
            license = licenses.gpl3;
            maintainers = [ ];
            mainProgram = "sfy";
            platforms = platforms.unix;
          };
        };
    in
    {
      # Packages
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = mkShellify pkgs;
          shellify = mkShellify pkgs;
        });

      # Overlay for use in NixOS/home-manager configurations
      overlays.default = final: prev: {
        shellify = mkShellify final;
      };

      # Home Manager module
      homeManagerModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.programs.shellify;
        in
        {
          options.programs.shellify = {
            enable = lib.mkEnableOption "shellify terminal multiplexer session manager";

            package = lib.mkOption {
              type = lib.types.package;
              default = self.packages.${pkgs.system}.default;
              defaultText = lib.literalExpression "inputs.shellify.packages.\${pkgs.system}.default";
              description = "The shellify package to use.";
            };
          };

          config = lib.mkIf cfg.enable {
            home.packages = [ cfg.package ];
          };
        };

      # Development shell
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              # Go toolchain (matches go.mod toolchain version)
              go_1_24
              goreleaser

              # Node for GUI (requires Node 22+)
              nodejs_22

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
              echo "  make lint         - Run linter (installs golangci-lint if needed)"
              echo "  make gui-dev      - Start GUI dev server"
            '';
          };
        });
    };
}
