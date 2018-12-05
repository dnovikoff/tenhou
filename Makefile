gobin:
	mkdir gobin

gobin/protoc-gen-go: gobin
	go build -o ./gobin/protoc-gen-go ./vendor/github.com/golang/protobuf/protoc-gen-go

# github.com/golang/protobuf/protoc-gen-go
CMD := protoc --plugin=protoc-gen-go=./gobin/protoc-gen-go --go_out=paths=source_relative,plugins=grpc:./genproto --proto_path=./proto

generate: gobin/protoc-gen-go
	rm -rf genproto
	mkdir genproto
	# cd proto &&
	$(CMD) ./proto/stats/*.proto

test:
	go test ./...
