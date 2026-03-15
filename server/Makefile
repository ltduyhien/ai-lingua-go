# := is the "immediately expanded" assignment; the value is set once when Make parses this line.
# Purpose: hold the path to our .proto sources so we can reuse it and change it in one place.
# Usage: referenced as $(PROTO_DIR) in the proto target.
PROTO_DIR := api/proto

# $(wildcard pattern) is a Make function that expands to a list of files matching the glob pattern.
# *.proto matches any file ending in .proto in the PROTO_DIR.
# Purpose: automatically pick up all proto files without listing them by name.
# Usage: passed to protoc so we generate from every proto we add.
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Same := assignment as above; stores the Go module path from go.mod.
# Purpose: protoc's module= option strips this prefix from go_package to produce output paths.
# Usage: without it, protoc would create github.com/ltduyhien/... under the output dir.
MODULE := github.com/ltduyhien/ai-lingua-go

# .PHONY marks targets that are not real files; Make always runs their recipe.
# If we omitted this and a file named "proto" existed, Make would skip the target thinking it is up to date.
# Purpose: ensure "make proto" always runs the generation.
.PHONY: proto

# "proto:" defines a target; the lines below (indented with tab) are the recipe run when you execute "make proto".
# Purpose: single command to regenerate Go code from .proto files.
proto:
	# Recipe lines run in the shell; \ continues the command onto the next line.
	# protoc: --go_out=. writes .pb.go to current dir; module= strips MODULE from go_package paths.
	# --go-grpc_out=. writes _grpc.pb.go; -I sets include path for proto imports; $(PROTO_FILES) are inputs.
	protoc --go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		-I $(PROTO_DIR) $(PROTO_FILES)
