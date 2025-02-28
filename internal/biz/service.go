package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ServiceRepoInterface interface {
	CommitWorkflow(context.Context, *Workflow) error
	GetWorkflow(context.Context, *Workflow) error
	ApplyService(context.Context, *Service, *ContinuousDeployment) error
	GetService(context.Context, *Service) error
}

type ServiceUseCase struct {
	repo ServiceRepoInterface
	log  *log.Helper
}

func NewServiceUseCase(repo ServiceRepoInterface, logger log.Logger) *ServiceUseCase {
	return &ServiceUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (w *Workflow) GetTask(taskName string) *WorkflowTask {
	for _, step := range w.WorkflowSteps {
		for _, task := range step.WorkflowTasks {
			if task.Name == taskName {
				return task
			}
		}
	}
	return nil
}

func (s *ServiceUseCase) CommitWorkflow(ctx context.Context, wf *Workflow) error {
	return s.repo.CommitWorkflow(ctx, wf)
}

func (s *ServiceUseCase) GetWorkflow(ctx context.Context, wf *Workflow) error {
	return s.repo.GetWorkflow(ctx, wf)
}

func (s *ServiceUseCase) ApplyService(ctx context.Context, service *Service, cd *ContinuousDeployment) error {
	return s.repo.ApplyService(ctx, service, cd)
}

func (s *ServiceUseCase) GetService(ctx context.Context, service *Service) error {
	return s.repo.GetService(ctx, service)
}
