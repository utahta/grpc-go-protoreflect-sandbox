package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/utahta/grpc-go-protoreflect-example/gen/option"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

const (
	address = "localhost:50050"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gc := grpcreflect.NewClient(ctx, rpb.NewServerReflectionClient(conn))

	services, err := gc.ListServices()
	if err != nil {
		log.Fatalf("failed to get services: %v", err)
	}

	enumDesc, err := tagEnumDescriptor()
	if err != nil {
		log.Fatalf("failed to get enum descriptor: %v", err)
	}
	log.Printf("enum: %v\n", enumDesc)

	for _, s := range services {
		if s == "grpc.reflection.v1alpha.ServerReflection" {
			continue
		}
		log.Printf("service: %v\n", s)

		fd, err := gc.FileContainingSymbol(s)
		if err != nil {
			log.Fatalf("failed to get file by symbol: %v", err)
		}
		log.Printf("package: %v\n", fd.GetPackage())
		log.Printf("services: %v\n", fd.GetServices())

		for _, srvDesc := range fd.GetServices() {
			for _, methodDesc := range srvDesc.GetMethods() {
				log.Printf("fullMethod: /%v.%v/%v\n", fd.GetPackage(), srvDesc.GetName(), methodDesc.GetName())

				switch {
				case proto.HasExtension(methodDesc.GetOptions(), option.E_Tag):
					tmp, err := proto.GetExtension(methodDesc.GetOptions(), option.E_Tag)
					if err != nil {
						log.Fatalf("failed to get ext: %v", err)
					}

					tag, ok := tmp.(*option.Tag)
					if !ok {
						log.Fatalf("failed to cast: %v", tmp)
					}
					log.Printf("tag: %v\n", tag)

				case proto.HasExtension(methodDesc.GetOptions(), option.E_Tags):
					tmp, err := proto.GetExtension(methodDesc.GetOptions(), option.E_Tags)
					if err != nil {
						log.Fatalf("failed to get ext: %v", err)
					}

					tags, ok := tmp.([]option.Tag)
					if !ok {
						log.Fatalf("failed to cast: %v", tmp)
					}
					log.Printf("tags: %v\n", tags)
				}
			}
		}
	}
}

func tagEnumDescriptor() (*descriptor.EnumDescriptorProto, error) {
	fdb, idxs := new(option.Tag).EnumDescriptor()
	if len(fdb) == 0 {
		return nil, fmt.Errorf("fdb is empty")
	}
	if len(idxs) == 0 {
		return nil, fmt.Errorf("idxs is empty")
	}

	r, err := gzip.NewReader(bytes.NewReader(fdb))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress: %v", err)
	}
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor: %v", err)
	}

	var fd descriptor.FileDescriptorProto
	if err := proto.Unmarshal(raw, &fd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal descriptor: %v", err)
	}

	idx := idxs[0]
	if len(fd.GetEnumType()) <= idx {
		return nil, fmt.Errorf("invalid enum type")
	}
	return fd.GetEnumType()[idx], nil
}
