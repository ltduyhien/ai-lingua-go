FROM golang:1.22-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
	protobuf-compiler \
	make \
	&& rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make proto && CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM debian:bookworm-slim

RUN useradd -D -s /sbin/nologin app

USER app

COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]
