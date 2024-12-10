package interfaces

import (
	"context"
	"errors"

	clusterApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/cluster"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type ClusterInterface struct {
	clusterApi.UnimplementedClusterInterfaceServer
	uc  *biz.ClusterUsecase
	log *log.Helper
}

func NewClusterInterface(uc *biz.ClusterUsecase, logger log.Logger) *ClusterInterface {
	return &ClusterInterface{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (c *ClusterInterface) CheckClusterInstalled(ctx context.Context, cluster *biz.Cluster) (*clusterApi.ClusterInstalled, error) {
	err := c.uc.CheckClusterInstalled(cluster)
	if errors.Is(err, biz.ErrClusterNotFound) {
		return &clusterApi.ClusterInstalled{
			Installed: false,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &clusterApi.ClusterInstalled{
		Installed: true,
	}, nil
}

func (c *ClusterInterface) CurrentCluster(ctx context.Context, cluster *biz.Cluster) (*biz.Cluster, error) {
	return c.uc.CurrentCluster(ctx, cluster)
}

func (c *ClusterInterface) HandlerNodes(ctx context.Context, cluster *biz.Cluster) (*biz.Cluster, error) {
	return c.uc.HandlerNodes(ctx, cluster)
}

func (c *ClusterInterface) MigrateToCluster(ctx context.Context, cluster *biz.Cluster) (*biz.Cluster, error) {
	return c.uc.MigrateToCluster(ctx, cluster)
}
