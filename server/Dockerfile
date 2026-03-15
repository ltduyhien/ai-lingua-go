FROM golang:1.23-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
	protobuf-compiler \
	make \
	&& rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make proto && go test ./... && CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM debian:bookworm-slim

RUN useradd -M -s /sbin/nologin app

USER app

COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]
