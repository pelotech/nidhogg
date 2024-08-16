# Build nidhogg manager binary

```
go build -o bin/manager github.com/uswitch/nidhogg/cmd/manager
```

# Execute unit tests and integration tests using kubebuilder

## Prerequisite : install kubebuilder

### Mac (arm64)

```
export KUBEBUILDER_VERSION=v4.1.1
export KUBEBUILDER_TOOLS_VERSION=1.30.0
export PLATFORM=darwin
export ARCH=arm64

sudo mkdir -p /usr/local/kubebuilder/bin
sudo curl -L "https://github.com/kubernetes-sigs/kubebuilder/releases/download/${KUBEBUILDER_VERSION}/kubebuilder_${PLATFORM}_${ARCH}" -o /usr/local/kubebuilder/bin/kubebuilder
sudo chmod +x /usr/local/kubebuilder/bin/kubebuilder

sudo curl -L "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-${KUBEBUILDER_TOOLS_VERSION}-${PLATFORM}-${ARCH}.tar.gz" | tar -xz -C /tmp/
sudo mv /tmp/kubebuilder/bin/* /usr/local/kubebuilder/bin && rm -rf /tmp/kubebuilder

# Put all the kubebuilder binaries in the allow list
sudo xattr -r -d com.apple.quarantine /usr/local/kubebuilder/bin
```

### Linux (amd64)

```
export KUBEBUILDER_VERSION=v4.1.1
export KUBEBUILDER_TOOLS_VERSION=1.30.0
export PLATFORM=linux
export ARCH=amd64

sudo mkdir -p /usr/local/kubebuilder/bin
sudo curl -L "https://github.com/kubernetes-sigs/kubebuilder/releases/download/${KUBEBUILDER_VERSION}/kubebuilder_${PLATFORM}_${ARCH}" -o /usr/local/kubebuilder/bin/kubebuilder
sudo chmod +x /usr/local/kubebuilder/bin/kubebuilder

sudo curl -L "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-${KUBEBUILDER_TOOLS_VERSION}-${PLATFORM}-${ARCH}.tar.gz" | tar -xz -C /tmp/
sudo mv /tmp/kubebuilder/bin/* /usr/local/kubebuilder/bin && rm -rf /tmp/kubebuilder
```

## Run tests

```
go test ./...
```