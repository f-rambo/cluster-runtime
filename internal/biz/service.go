package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ServiceRepoInterface interface {
	CommitWorkflow(context.Context, *Workflow) error
	GetWorkflow(context.Context, *Workflow) error
	CleanWorkflow(context.Context, *Workflow) error
	ApplyService(context.Context, *Service, *ContinuousDeployment) error
	GetService(context.Context, *Service) error
}

type ServiceUseCase struct {
	repo ServiceRepoInterface
	log  *log.Helper
}

func (w *Workflow) GetFirstStep() *WorkflowStep {
	if len(w.WorkflowSteps) == 0 {
		return nil
	}
	var firstStep *WorkflowStep
	for _, step := range w.WorkflowSteps {
		if firstStep == nil || step.Order < firstStep.Order {
			firstStep = step
		}
	}
	return firstStep
}

func (w *Workflow) GetNextStep(s *WorkflowStep) *WorkflowStep {
	if len(w.WorkflowSteps) == 0 {
		return nil
	}
	var nextStep *WorkflowStep
	for _, step := range w.WorkflowSteps {
		if step.Order > s.Order {
			if nextStep == nil || step.Order < nextStep.Order {
				nextStep = step
			}
		}
	}
	return nextStep
}

func (s *WorkflowStep) GetFirstTask() *WorkflowTask {
	if len(s.WorkflowTasks) == 0 {
		return nil
	}
	var firstTask *WorkflowTask
	for _, task := range s.WorkflowTasks {
		if firstTask == nil || task.Order < firstTask.Order {
			firstTask = task
		}
	}
	return firstTask
}

func (s *WorkflowStep) GetNextTask(t *WorkflowTask) *WorkflowTask {
	if len(s.WorkflowTasks) == 0 {
		return nil
	}
	var nextTask *WorkflowTask
	for _, task := range s.WorkflowTasks {
		if task.Order > t.Order {
			if nextTask == nil || task.Order < nextTask.Order {
				nextTask = task
			}
		}
	}
	return nextTask
}

func (w *Workflow) GetStorageName() string {
	return w.Name + "-storage"
}

func (w *Workflow) GetWorkdir() string {
	return "/app"
}

func (w *Workflow) GetWorkdirName() string {
	return "app"
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

func (s *ServiceUseCase) CleanWorkflow(ctx context.Context, wf *Workflow) error {
	return s.repo.CleanWorkflow(ctx, wf)
}

func (s *ServiceUseCase) ApplyService(ctx context.Context, service *Service, cd *ContinuousDeployment) error {
	return s.repo.ApplyService(ctx, service, cd)
}

func (s *ServiceUseCase) GetService(ctx context.Context, service *Service) error {
	return s.repo.GetService(ctx, service)
}
