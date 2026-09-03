# Build the manager binary.
#
# The builder stage always runs on the build host's platform and cross-compiles
# for the target (CGO is off, so the Go toolchain needs no target sysroot). A
# multi-platform buildx invocation therefore compiles natively for every
# target instead of emulating each one.
FROM --platform=$BUILDPLATFORM golang:1.26 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Module manifests first so the dependency download layer survives source
# changes. go.work pins the operator as its own workspace (see that file).
COPY go.mod go.sum go.work ./
RUN go mod download

# Go source (filtered by .dockerignore)
COPY . .

# GOARCH intentionally has no default: with a plain `docker build`, buildx sets
# TARGETARCH to the host's architecture, so the binary always matches the
# image platform it ships in.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
