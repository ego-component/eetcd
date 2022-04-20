package main

import (
	"context"

	"github.com/ego-component/eetcd"
	"github.com/ego-component/eetcd/examples/helloworld"
	"github.com/ego-component/eetcd/registry"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/server/egrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().
		Registry(registry.Load("registry").Build(registry.WithClientEtcd(eetcd.Load("etcd").Build()))).
		Serve(func() server.Server {
			component := egrpc.Load("server.grpc").Build()
			helloworld.RegisterGreeterServer(component.Server, &Greeter{})
			return component
		}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	if request.Name == "error" {
		return nil, status.Error(codes.Unavailable, "error")
	}

	return &helloworld.HelloResponse{
		Message: "Hello EGO",
	}, nil
}
