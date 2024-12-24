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
            golangci-lint.enable = true;
            gofmt.enable = true;
            gotest.enable = true;
          };
        };
      in {
      # devShells."x86_64-linux".default = import ./shell.nix { inherit pkgs; };
      devShells.default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          inputs.mockery.legacyPackages.${system}.go-mockery
          go_1_23
          air
          go-swag
          goose sqlc
          delve
          golangci-lint

          docker
          docker-compose
          gnumake
          pkg-config
          git
          geos
          proj
          gitflow
        ];

      shellHook = ''
        export PATH="$(go env GOPATH)/bin:$PATH"
        go install github.com/jmattheis/goverter/cmd/goverter@latest

        go mod download
        git flow init -d -f -t v

        ${pre-commit-check.shellHook}
      '';
      };
    }
  ));
}
