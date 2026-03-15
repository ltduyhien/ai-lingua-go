PROTO_DIR := api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
MODULE := github.com/ltduyhien/ai-lingua-go

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		-I $(PROTO_DIR) $(PROTO_FILES)
