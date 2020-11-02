.PHONY: generate
all: generate test build

include protoc.mk

gobin:
	mkdir gobin

# github.com/golang/protobuf/protoc-gen-go
CMD := $(protoc_go_cmd) --go_out=paths=source_relative,plugins=grpc:./genproto --proto_path=./proto

.PHONY: generate
generate: $(protoc_gen_go)
	rm -rf genproto
	mkdir genproto
	# cd proto &&
	$(CMD) ./proto/stats/*.proto

.PHONY: test
test:
	go test -mod vendor ./...

.PHONY: build
build:
	mkdir -p build
	GOBIN=$(shell pwd)/build go install -mod vendor ./cmd/...

.PHONY: testcover
testcover:
	go test -mod vendor -race -coverprofile=coverage.txt -covermode=atomic ./...
