# Build the manager binary
FROM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e as builder

# Copy in the go src
WORKDIR /app

COPY pkg/    pkg/
COPY cmd/    cmd/
COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o manager github.com/uswitch/nidhogg/cmd/manager

# Copy the controller-manager into a thin image
FROM gcr.io/distroless/static:nonroot@sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3
WORKDIR /
COPY --from=builder /app/manager .
ENTRYPOINT ["/manager"]
