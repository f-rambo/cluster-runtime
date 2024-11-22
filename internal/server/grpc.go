package server

import (
	"time"

	appApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/app"
	cluster "github.com/f-rambo/cloud-copilot/cluster-runtime/api/cluster"
	logApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/log"
	projectApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/project"
	serviceApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/service"
	userApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/user"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/conf"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/interfaces"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	logInterface *interfaces.LogInterface,
	clusterInterface *interfaces.ClusterInterface,
	appInterface *interfaces.AppInterface,
	projectInterface *interfaces.ProjectInterface,
	serviceInterface *interfaces.ServiceInterface,
	userInterface *interfaces.UserInterface,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.GRPC.Network != "" {
		opts = append(opts, grpc.Network(c.GRPC.Network))
	}
	if c.GRPC.Addr != "" {
		opts = append(opts, grpc.Address(c.GRPC.Addr))
	}
	if c.GRPC.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(c.GRPC.Timeout)*time.Second))
	}
	srv := grpc.NewServer(opts...)
	cluster.RegisterClusterInterfaceServer(srv, clusterInterface)
	logApi.RegisterLogInterfaceServer(srv, logInterface)
	appApi.RegisterAppInterfaceServer(srv, appInterface)
	projectApi.RegisterProjectInterfaceServer(srv, projectInterface)
	serviceApi.RegisterServiceInterfaceServer(srv, serviceInterface)
	userApi.RegisterUserInterfaceServer(srv, userInterface)
	return srv
}
