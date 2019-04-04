dep:
	dep ensure -v

.PHONY: proto
proto:
	rm -rf gen && mkdir gen
	protoc -I./proto --go_out=plugins=grpc:${GOPATH}/src proto/helloworld/*.proto
	protoc -I./proto --go_out=plugins=grpc:${GOPATH}/src proto/option/*.proto

serv:
	go run ./server/main.go

