all: gen-go-proto gen-cypher

gen-go-proto:
	protoc --proto_path=./proto --go_out=plugins=grpc:./service service.proto

# js-proto:
# 	mkdir -p ./service/js
# 	protoc --proto_path=./service --js_out=import_style=commonjs:./service/js service.proto

build:
	go build -o build/draft-server cmd/server/*.go

gen-cypher:
	pigeon graph/parser/cypher/cypher.peg | goimports > graph/parser/cypher/cypher.go

fmt:
	gofmt -w ./..

test: gen-go-proto gen-cypher
	go test -timeout 30s -race -v ./... -cover