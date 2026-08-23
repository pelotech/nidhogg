# Build the manager binary
FROM golang:1.26.6@sha256:0d1d3a794be25f809dd2cb3160d8c73276c4056a9f8242a138e908ddeee7b6b6 as builder

# Copy in the go src
WORKDIR /app

COPY pkg/    pkg/
COPY cmd/    cmd/
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o manager github.com/uswitch/nidhogg/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest@sha256:678c6550cc43645e08669028bc177f50be4e7c5b8cca677067b1914d4afc7a03
WORKDIR /
COPY --from=builder /app/manager .
ENTRYPOINT ["/manager"]
