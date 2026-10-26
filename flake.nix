{
  description = "nidhogg dev shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [ inputs.treefmt-nix.flakeModule ];
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      perSystem =
        { pkgs, ... }:
        {
          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              # go
              go
              setup-envtest

              # helm
              kubernetes-helm
              helm-docs

              # tasks
              just

              # pre-commit
              prek
              yamllint
              actionlint
              zizmor
            ];

            # envtest needs kube-apiserver and etcd; nixpkgs has no darwin
            # kube-apiserver, so fetch the upstream envtest binaries instead.
            shellHook = ''
              export KUBEBUILDER_ASSETS="$(setup-envtest use -p path 1.37.x)"
            '';
          };

          treefmt = {
            projectRootFile = "flake.nix";
            programs.nixfmt.enable = true;
          };
        };
    };
}
