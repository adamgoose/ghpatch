{ inputs, cell }:
let
  # The `inputs` attribute allows us to access all of our flake inputs.
  inherit (inputs) nixpkgs std;

  # This is a common idiom for combining lib with builtins.
  l = nixpkgs.lib // builtins;
in
rec {
  default = ghpatch;

  ghpatch = with cell.pkgs.default; buildGoApplication rec {
    pname = "ghpatch";
    version = "0.0.1";
    pwd = inputs.self;
    src = inputs.self;
    modules = "${inputs.self}/gomod2nix.toml";
    doCheck = false;
  };
}
