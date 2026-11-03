{
  description = "Golang Template";

  inputs = { nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-26.05"; };

  outputs = { self, nixpkgs }:
    let
      # define system once
      system = "x86_64-linux";
      # use it here, and bind platform-specific packages to `pkgs`
      pkgs = nixpkgs.legacyPackages.${system};
      nativeBuildInputs = with pkgs; [ go ];
      developmentPakcages = with pkgs; [ gopls delve ];
    in {

      overlays.default = final: prev: {
        torrentHelperBot = pkgs.buildGoModule  {
          pname = "Torrent Helper Bot";
          version = "0.0.1";
          src = ./.;
          vendorHash = "sha256-iGtV8n1mo7OBtdHgnLs5jqCwFLoaKYVLIZmYAiuVxH0=";
          env = { CGO_ENABLED = 0; };
        };
      };
      packages.${system}.torrentHelperBot =
        (import nixpkgs { inherit system; overlays = [ self.overlays.default ]; }).torrentHelperBot;

      devShells.${system}.default = pkgs.mkShell {
        packages = nativeBuildInputs ++ developmentPakcages
          ++ [ pkgs.bashInteractive ];
      };
      nixosModules.torrentHelperBot = import ./nixosModule;
    };
}
