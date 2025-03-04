package repo

import (
	"context"

	argoworkflow "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	coreV1 "k8s.io/api/core/v1"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	mateV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MiB = 1024 * 1024 // 1 MiB = 2^20 字节

const GI = 1024 * 1024 * 1024 // 1 Gi = 2^30 字节

type ServiceRepo struct {
	log *log.Helper
}

func NewServiceRepo(logger log.Logger) biz.ServiceRepoInterface {
	return &ServiceRepo{
		log: log.NewHelper(logger),
	}
}

func (s *ServiceRepo) CommitWorkflow(ctx context.Context, wf *biz.Workflow) error {
	argoClient, err := NewArgoWorkflowClient(wf.Namespace)
	if err != nil {
		return err
	}
	argoWf, err := argoClient.Create(ctx, ConvertToArgoWorkflow(wf))
	if err != nil {
		return err
	}
	kubeClient, err := GetKubernetesClientSet()
	if err != nil {
		return err
	}
	pvc := &coreV1.PersistentVolumeClaim{
		ObjectMeta: mateV1.ObjectMeta{
			Name:      wf.GetStorageName(),
			Namespace: wf.Namespace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			AccessModes: []coreV1.PersistentVolumeAccessMode{
				coreV1.ReadWriteOnce,
			},
			Resources: coreV1.VolumeResourceRequirements{
				Requests: coreV1.ResourceList{
					coreV1.ResourceStorage: *resource.NewQuantity(int64(1)*GI, resource.BinarySI),
				},
			},
			StorageClassName: &wf.StorageClass,
		},
	}
	pvc.OwnerReferences = []mateV1.OwnerReference{
		{
			APIVersion:         argoworkflow.APIVersion,
			Kind:               argoworkflow.WorkflowKind,
			Name:               argoWf.Name,
			UID:                argoWf.UID,
			Controller:         utils.PointerBool(true),
			BlockOwnerDeletion: utils.PointerBool(true),
		},
	}
	_, err = kubeClient.CoreV1().PersistentVolumeClaims(wf.Namespace).Create(ctx, pvc, mateV1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceRepo) GetWorkflow(ctx context.Context, wf *biz.Workflow) error {
	argoClient, err := NewArgoWorkflowClient(wf.Namespace)
	if err != nil {
		return err
	}
	argoWf, err := argoClient.Get(ctx, wf.Name, mateV1.GetOptions{})
	if err != nil {
		return err
	}
	SetWorkflowStatus(argoWf, wf)
	return nil
}

func (s *ServiceRepo) CleanWorkflow(ctx context.Context, wf *biz.Workflow) error {
	argoClient, err := NewArgoWorkflowClient(wf.Namespace)
	if err != nil {
		return err
	}
	err = argoClient.Delete(ctx, wf.Name)
	if err != nil && !k8sError.IsNotFound(err) {
		return err
	}
	kubeClient, err := GetKubernetesClientSet()
	if err != nil {
		return err
	}
	err = kubeClient.CoreV1().PersistentVolumeClaims(wf.Namespace).Delete(ctx, wf.GetStorageName(), mateV1.DeleteOptions{})
	if err != nil && !k8sError.IsNotFound(err) {
		return err
	}
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
