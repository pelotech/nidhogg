# Dev environment

Tooling (go, helm, envtest, prek, just, ...) comes from the nix flake. With
[nix](https://nixos.org) and [direnv](https://direnv.net) installed:

```
direnv allow
prek install --hook-type pre-commit --hook-type commit-msg
```

Without direnv, run `nix develop` instead. Entering the shell sets
`KUBEBUILDER_ASSETS` to the envtest kube-apiserver/etcd binaries.

# Build nidhogg manager binary

```
just build
```

# Run unit and integration tests

```
just test
```

# Run pre-commit hooks

```
just lint
```
