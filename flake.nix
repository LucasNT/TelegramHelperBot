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

      packages.${system}.default = pkgs.buildGoModule {
        pname = "Torrent Helper Bot";
        version = "0.0.1";
        src = ./.;
        vendorHash = "sha256-iGtV8n1mo7OBtdHgnLs5jqCwFLoaKYVLIZmYAiuVxH0=";
        # subPackages = [ "path1", "path2" ]
        env = { CGO_ENABLED = 0; };
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = nativeBuildInputs ++ developmentPakcages
          ++ [ pkgs.bashInteractive ];
      };
      nixosModules = import ./nixosModule;
    };
}
