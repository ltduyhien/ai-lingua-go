# Stage 1: build the binary (proto generation + Go build).
# We use the official Go image so we have the Go toolchain; we add protoc and Make for proto generation.
FROM golang:1.22-bookworm AS builder

# Install protoc (protocol buffer compiler) and make from Debian packages.
# apt-get update refreshes package lists; -y avoids interactive prompts.
# protobuf-compiler provides the protoc binary; make runs our Makefile.
RUN apt-get update && apt-get install -y --no-install-recommends \
	protobuf-compiler \
	make \
	&& rm -rf /var/lib/apt/lists/*

# Install the Go plugins that protoc uses to generate Go code.
# They are installed into $GOPATH/bin (or $HOME/go/bin) so protoc can find them.
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set the working directory inside the container to /app so all following paths are relative to it.
WORKDIR /app

# Copy go.mod and go.sum first so we can download dependencies in a separate layer.
# Docker caches layers; if only source changes, we reuse the dependency layer.
COPY go.mod go.sum ./

# Download Go module dependencies; they are stored in the module cache inside the image.
RUN go mod download

# Copy the rest of the source (proto files, Makefile, config, cmd, internal, etc.).
COPY . .

# Generate Go code from .proto files; creates api/gen/ with .pb.go and _grpc.pb.go.
# Then build the server binary into /app/server; CGO_ENABLED=0 produces a static binary for the final image.
RUN make proto && CGO_ENABLED=0 go build -o /app/server ./cmd/server

# Stage 2: minimal runtime image; we only need the binary and no build tools.
# scratch is an empty image; we add only the binary and (optional) CA certs if we need HTTPS.
FROM debian:bookworm-slim

# Create a non-root user so the process does not run as root.
# -D means no password, no home dir; the name of the user is "app".
RUN useradd -D -s /sbin/nologin app

# Switch to the app user so the binary runs with reduced privileges.
USER app

# Copy the server binary from the builder stage into this image.
# Only this file is needed to run the service; no Go toolchain or source.
COPY --from=builder /app/server /server

# Default command when the container starts; runs the gRPC server.
# Override with docker run ... or in docker-compose to pass flags or env.
ENTRYPOINT ["/server"]
