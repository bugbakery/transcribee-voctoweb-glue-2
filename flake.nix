{
  description = "transcribee voctoweb glue";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      self,
      ...
    }:
    {
      overlays.default = (
        final: prev:
        let
          pkgs = import nixpkgs {
            system = final.system;
          };
          lib = nixpkgs.lib;
        in
        {
          transcribee-voctoweb = import ./nix/pkgs/backend.nix {
            inherit
              pkgs
              lib
              ;
          };
        }
      );

      nixosModules.default = {
        nixpkgs.overlays = [ self.overlays.default ];
      };
    }
    // (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      in
      {
        packages = {
          transcribee-voctoweb = pkgs.transcribee-voctoweb;
        };

        formatter = pkgs.nixfmt-rfc-style;

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gcc

            pkgs.nodejs_22

            # nix tooling
            pkgs.nixpkgs-fmt

            # dev tooling
            pkgs.overmind
          ];
        };
      }
    ));
}
