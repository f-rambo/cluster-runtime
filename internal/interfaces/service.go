package interfaces

import (
	"context"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/api/common"
	serviceApi "github.com/f-rambo/cloud-copilot/cluster-runtime/api/service"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
)

type ServiceInterface struct {
	serviceApi.UnimplementedServiceInterfaceServer
	serviceUc *biz.ServiceUseCase
}

func NewServiceInterface(serviceUc *biz.ServiceUseCase) *ServiceInterface {
	return &ServiceInterface{
		serviceUc: serviceUc,
	}
}

func (s *ServiceInterface) Create(ctx context.Context, req *serviceApi.CreateReq) (*common.Msg, error) {
	err := s.serviceUc.Create(ctx, req.Namespace, req.Workflow)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (s *ServiceInterface) GenerateCIWorkflow(ctx context.Context, service *biz.Service) (*serviceApi.GenerateCIWorkflowResponse, error) {
	ciWf, cdwf, err := s.serviceUc.GenerateCIWorkflow(ctx, service)
	if err != nil {
		return nil, err
	}
	return &serviceApi.GenerateCIWorkflowResponse{
		CiWorkflow: ciWf,
		CdWorkflow: cdwf,
	}, nil
}
