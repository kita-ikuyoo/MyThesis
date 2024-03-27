package main

import (
	"IHome/GetImageCd/handler"
	pb "IHome/GetImageCd/proto"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4/registry"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	service = "go.micro.server.GetImageCd"
	version = "latest"
)

func main() {
	// Create service
	consulRegistry := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"127.0.0.1:8500", // 这里假设您的Consul服务运行在本机的8500端口
		}
	})
	server := micro.NewService(

		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()), //
		micro.Registry(consulRegistry),
		micro.Name(service),
	)

	server.Init()

	// Register handler
	if err := pb.RegisterGetImageCdHandler(server.Server(), new(handler.GetImageCd)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := server.Run(); err != nil {
		logger.Fatal(err)
	}
}
