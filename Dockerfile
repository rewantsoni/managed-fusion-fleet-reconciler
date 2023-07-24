# Build stage 1
FROM golang:1.20 as builder

WORKDIR /workspace

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/managed-fusion-fleet-reconciler main.go

# Build stage 2

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /workspace/bin/managed-fusion-fleet-reconciler /usr/local/bin/managed-fusion-fleet-reconciler

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/managed-fusion-fleet-reconciler"]
