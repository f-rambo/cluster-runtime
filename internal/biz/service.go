package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/go-kratos/kratos/v2/log"
)

type ServiceUseCase struct {
	log *log.Helper
}

func NewServiceUseCase(logger log.Logger) *ServiceUseCase {
	return &ServiceUseCase{
		log: log.NewHelper(logger),
	}
}

func (s *ServiceUseCase) GenerateCIWorkflow(ctx context.Context, service *Service) (ciWf *Workflow, cdwf *Workflow, err error) {
	ciWorkflowJson, err := GetDefaultWorklfows(ctx, strings.ToLower(service.Business), strings.ToLower(service.Technology), WorkflowType_ContinuousIntegration.String())
	if err != nil {
		return nil, nil, err
	}
	ciWf = &Workflow{
		Name:     fmt.Sprintf("default-%s-%s-%s", service.Business, service.Technology, WorkflowType_ContinuousIntegration.String()),
		Workflow: ciWorkflowJson,
	}
	cdWorkflowJson, err := GetDefaultWorklfows(ctx, strings.ToLower(service.Business), strings.ToLower(service.Technology), WorkflowType_ContinuousDeployment.String())
	if err != nil {
		return nil, nil, err
	}
	cdwf = &Workflow{
		Name:     fmt.Sprintf("default-%s-%s-%s", service.Business, service.Technology, WorkflowType_ContinuousDeployment.String()),
		Workflow: cdWorkflowJson,
	}
	return ciWf, cdwf, nil
}

func (s *ServiceUseCase) Create(ctx context.Context, namespace string, workflow *Workflow) error {
	kubeConf, err := getKubeConfig()
	if err != nil {
		return err
	}
	argoClient, err := newForConfig(kubeConf)
	if err != nil {
		return err
	}
	argoWf := &wfv1.Workflow{}
	err = json.Unmarshal(workflow.Workflow, argoWf)
	if err != nil {
		return err
	}
	_, err = argoClient.Workflows(namespace).Create(ctx, argoWf)
	if err != nil {
		return err
	}
	return nil
}
