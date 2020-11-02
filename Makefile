include protoc.mk

gobin:
	mkdir gobin

# github.com/golang/protobuf/protoc-gen-go
CMD := $(protoc_go_cmd) --go_out=paths=source_relative,plugins=grpc:./genproto --proto_path=./proto

generate: $(protoc_gen_go)
	rm -rf genproto
	mkdir genproto
	# cd proto &&
	$(CMD) ./proto/stats/*.proto

test:
	go test ./...
