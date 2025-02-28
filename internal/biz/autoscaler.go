package biz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

func (c *Cluster) generateNodeLables(nodeGroup *NodeGroup) string {
	lableMap := make(map[string]string)
	lableMap["cluster"] = c.Name
	lableMap["cluster_id"] = cast.ToString(c.Id)
	lableMap["cluster_type"] = c.Type.String()
	lableMap["region"] = c.Region
	lableMap["nodegroup"] = nodeGroup.Name
	lableMap["nodegroup_type"] = nodeGroup.Type.String()
	lablebytes, _ := json.Marshal(lableMap)
	return string(lablebytes)
}

func (uc *ClusterUsecase) GetCurrentCluster(ctx context.Context) (cluster *Cluster, err error) {
	err = uc.CurrentCluster(ctx, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (uc *ClusterUsecase) NodeGroupIncreaseSize(ctx context.Context, cluster *Cluster, nodeGroup *NodeGroup, size int32) error {
	for i := 0; i < int(size); i++ {
		node := &Node{
			Name:        fmt.Sprintf("%s-%s", cluster.Name, uuid.New().String()),
			Role:        NodeRole_WORKER,
			Status:      NodeStatus_NODE_CREATING,
			ClusterId:   cluster.Id,
			NodeGroupId: nodeGroup.Id,
		}
		cluster.Nodes = append(cluster.Nodes, node)
	}
	return nil
}

func (uc *ClusterUsecase) DeleteNodes(ctx context.Context, cluster *Cluster, nodes []*Node) error {
	for _, node := range nodes {
		for i, n := range cluster.Nodes {
			if n.Id == node.Id {
				cluster.Nodes = append(cluster.Nodes[:i], cluster.Nodes[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (uc *ClusterUsecase) NodeGroupTemplateNodeInfo(ctx context.Context, cluster *Cluster, nodeGroup *NodeGroup) (*Node, error) {
	return &Node{
		Name:        fmt.Sprintf("%s-%s", cluster.Name, uuid.New().String()),
		Role:        NodeRole_WORKER,
		Status:      NodeStatus_NODE_CREATING,
		ClusterId:   cluster.Id,
		NodeGroupId: nodeGroup.Id,
		Labels:      cluster.generateNodeLables(nodeGroup),
	}, nil
}

func (uc *ClusterUsecase) Cleanup(ctx context.Context) error {
	return nil
}

func (uc *ClusterUsecase) Refresh(ctx context.Context) (err error) {
	cluster := &Cluster{}
	err = uc.CurrentCluster(ctx, cluster)
	if err != nil {
		return err
	}
	// rpc to cloud cluster
	return nil
}
