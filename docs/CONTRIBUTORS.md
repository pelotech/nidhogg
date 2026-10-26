# Dev environment

Tooling (go, helm, envtest, prek, just, ...) comes from the nix flake. With
[nix](https://nixos.org) and [direnv](https://direnv.net) installed:

```
direnv allow
prek install --hook-type pre-commit --hook-type commit-msg
```

Without direnv, run `nix develop` instead.

# Build nidhogg manager binary

```
just build
```

# Run unit and integration tests

`just test` downloads the envtest kube-apiserver/etcd binaries via
`setup-envtest` (nixpkgs has no darwin kube-apiserver) and points
`KUBEBUILDER_ASSETS` at them.

```
just test
```

# Run pre-commit hooks

```
just lint
```
