package repo

import (
	"context"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	mateV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceRepo struct {
	log *log.Helper
}

func NewServiceRepo(logger log.Logger) biz.ServiceRepoInterface {
	return &ServiceRepo{
		log: log.NewHelper(logger),
	}
}

func (s *ServiceRepo) CommitWorklfow(ctx context.Context, wf *biz.Workflow) error {
	argoClient, err := NewArgoWorkflowClient(wf.Namespace)
	if err != nil {
		return err
	}
	argoWf, err := argoClient.Create(ctx, ConvertToArgoWorkflow(wf))
	if err != nil {
		return err
	}
	wf.RuntimeName = argoWf.GenerateName
	wf.Lables = utils.MapToString(argoWf.Labels)
	return nil
}

func (s *ServiceRepo) GetWorkflow(ctx context.Context, wf *biz.Workflow) error {
	argoClient, err := NewArgoWorkflowClient(wf.Namespace)
	if err != nil {
		return err
	}
	argoWf, err := argoClient.Get(ctx, wf.RuntimeName, mateV1.GetOptions{})
	if err != nil {
		return err
	}
	SetWorkflowStatus(argoWf, wf)
	return nil
}

func (s *ServiceRepo) ApplyService(ctx context.Context, service *biz.Service, cd *biz.ContinuousDeployment) error {
	cloudServiceClinet, err := NewClServiceClient(service.Namespace)
	if err != nil {
		return err
	}
	cloudService, err := cloudServiceClinet.Get(ctx, service.Name, mateV1.GetOptions{})
	if err != nil {
		if k8sError.IsNotFound(err) {
			_, err := cloudServiceClinet.Create(ctx, ConvertToCloudService(service, cd))
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	cs := ConvertToCloudService(service, cd)
	cs.ResourceVersion = cloudService.ResourceVersion
	_, err = cloudServiceClinet.Update(ctx, cs)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceRepo) GetService(ctx context.Context, service *biz.Service) error {
	cloudServiceClinet, err := NewClServiceClient(service.Namespace)
	if err != nil {
		return err
	}
	cloudService, err := cloudServiceClinet.Get(ctx, service.Name, mateV1.GetOptions{})
	if err != nil {
		return err
	}
	SetServiceStatus(service, cloudService)
	return nil
}
