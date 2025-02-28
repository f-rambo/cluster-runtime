//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/conf"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/interfaces"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/repo"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/server"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(repo.ProviderSet, server.ProviderSet, biz.ProviderSet, interfaces.ProviderSet, newApp))
}
