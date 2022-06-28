all: build test

build: proto build-server build-client build-kong-plugin

proto:
	protoc \
		-I. \
		-I${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.7/ \
		--go_out=. \
		--go_opt=paths=source_relative \
    	--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
    ./example/person/person.proto

build-server:
	go build -o hello-server ./example/grpc/server/

build-client:
	go build -o hello-client ./example/grpc/client/

build-kong-plugin:
	go build -o kong-plugin-validate-protobuf ./plugin/kong/

test:
	go test ./...
