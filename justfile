# List recipes
default:
    @just --list

# Build the manager binary
build:
    go build -o bin/manager github.com/uswitch/nidhogg/cmd/manager

# Run unit and envtest integration tests
test:
    go test ./pkg/... ./cmd/...

# Run all pre-commit hooks
lint:
    prek run --all-files
