package interfaces

import (
	"context"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/api/common"
	projectApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/project"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProjectInterface struct {
	projectApi.UnimplementedProjectInterfaceServer
	projectUc *biz.ProjectUseCase
	log       *log.Helper
}

func NewProjectInterface(uc *biz.ProjectUseCase, logger log.Logger) *ProjectInterface {
	return &ProjectInterface{
		projectUc: uc,
		log:       log.NewHelper(logger),
	}
}

func (p *ProjectInterface) CreateNamespace(ctx context.Context, req *projectApi.CreateNamespaceReq) (*common.Msg, error) {
	if req.Namespace == "" {
		return nil, nil
	}
	err := p.projectUc.CreateNamespace(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (p *ProjectInterface) GetNamespaces(ctx context.Context, _ *emptypb.Empty) (*projectApi.Namesapces, error) {
	namespaces, err := p.projectUc.GetNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	return &projectApi.Namesapces{Namespaces: namespaces}, nil
}
