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

func (s *ServiceInterface) ApplyService(ctx context.Context, req *serviceApi.ApplyServiceRequest) (*common.Msg, error) {
	err := s.serviceUc.ApplyService(ctx, req.Service, req.Cd)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (s *ServiceInterface) GetService(ctx context.Context, service *biz.Service) (*biz.Service, error) {
	err := s.serviceUc.GetService(ctx, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceInterface) CommitWorklfow(ctx context.Context, wf *biz.Workflow) (*common.Msg, error) {
	err := s.serviceUc.CommitWorklfow(ctx, wf)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

func (s *ServiceInterface) GetWorkflow(ctx context.Context, wf *biz.Workflow) (*biz.Workflow, error) {
	err := s.serviceUc.GetWorkflow(ctx, wf)
	return wf, err
}
