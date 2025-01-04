{
  description = "Development enviroment for green ecolution backend";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    mockery.url = "github:nixos/nixpkgs/c3392ad349a5227f4a3464dce87bcc5046692fce";
    pre-commit-hooks.url = "github:cachix/git-hooks.nix";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs: 
    (flake-utils.lib.eachDefaultSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
        pre-commit-check = inputs.pre-commit-hooks.lib.${system}.run {
          src = ./.;
          hooks = {
            golangci-lint.enable = false;
            gofmt.enable = true;
            gotest.enable = false;
          };
        };

        goverter = pkgs.buildGoModule rec {
          pname = "goverter";
          version = "1.7.0";

          src = pkgs.fetchFromGitHub {
            owner = "jmattheis";
            repo = "goverter";
            tag = "v${version}";
            sha256 = "sha256-VgwmnB6FP7hlUrZpKun38T4K2YSDl9yYuMjdzsEhCF4=";
          };

          vendorHash = "sha256-uQ1qKZLRwsgXKqSAERSqf+1cYKp6MTeVbfGs+qcdakE=";
          subPackages = [ "cmd/goverter" ];
        };
      in {
      # devShells."x86_64-linux".default = import ./shell.nix { inherit pkgs; };
      devShells.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          inputs.mockery.legacyPackages.${system}.go-mockery
          go_1_23
          air
          go-swag
          goose
          sqlc
          delve
          golangci-lint
          goverter

          docker
          docker-compose
          gnumake
          pkg-config
          git
          geos
          proj
          gitflow
          yq-go
        ];

      shellHook = ''
        go mod download
        ${pre-commit-check.shellHook}
      '';
      };
    }
  ));
}
