{
  pkgs,
  ...
}:
pkgs.buildGoModule {
  pname = "transcribee-voctoweb";
  version = "0.1.0";
  src = ../..;
  vendorHash = "sha256-HfkWQw6bU5UglgzTmSiclYl6yaacjt6qFyoDgdMDH/o=";
}
