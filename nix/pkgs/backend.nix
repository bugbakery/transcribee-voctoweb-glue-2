{
  pkgs,
  lib,
  ...
}:
let
  frontend = pkgs.buildNpmPackage (
    lib.fix (self: {
      pname = "transcribee-voctoweb-frontend";
      version = "0.1.0";
      src = ../../frontend;

      nativeBuildInputs = [ pkgs.git ];

      npmDeps = pkgs.importNpmLock {
        npmRoot = self.src;
      };

      installPhase = ''
        runHook preInstall
        cp -r dist $out
        runHook postInstall
      '';

      npmBuildScript = "build";

      npmConfigHook = pkgs.importNpmLock.npmConfigHook;
    })
  );
in
pkgs.buildGoModule {
  pname = "transcribee-voctoweb";
  version = "0.1.0";
  src = ../..;
  vendorHash = "sha256-HfkWQw6bU5UglgzTmSiclYl6yaacjt6qFyoDgdMDH/o=";

  postInstall = ''
    cp -R ${frontend} $out/pb_public
'';
}
