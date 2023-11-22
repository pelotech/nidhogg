# Build the manager binary
FROM golang:1.21.4 as builder

# Copy in the go src
WORKDIR /app

COPY pkg/    pkg/
COPY cmd/    cmd/
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/uswitch/nidhogg/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /app/manager .
ENTRYPOINT ["/manager"]
