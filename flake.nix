{
  description = "Winecellar flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-24.11";
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in {
    formatter.x86_64-linux = pkgs.alejandra;
    packages.x86_64-linux.winecellar = with pkgs;
      buildGoModule rec {
        pname = "winecellar";
        version = "0.0.1";
        src = ./.;
        vendorHash = null;
      };
    packages.x86_64-linux.default = self.packages.x86_64-linux.winecellar;
    devShells.x86_64-linux.default = with pkgs;
      mkShellNoCC {
        buildInputs = with pkgs; [go];
        packages = with pkgs; [gopls gotools golangci-lint sqlite sqlc goose];
        shellHook = ''echo "Welcome to the winecellar development environment"'';
        GOOSE_DRIVER = "sqlite3";
        GOOSE_MIGRATION_DIR = "./sql/schema/";
      };
  };
}
