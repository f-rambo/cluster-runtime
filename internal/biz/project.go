package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectUseCase struct {
	log *log.Helper
}

func NewProjectUseCase(logger log.Logger) *ProjectUseCase {
	return &ProjectUseCase{
		log: log.NewHelper(logger),
	}
}

func (p *ProjectUseCase) CreateNamespace(ctx context.Context, namespace string) error {
	kubeClientSet, err := GetKubeClientByInCluster()
	if err != nil {
		return err
	}
	_, err = kubeClientSet.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})
	return err
}

func (p *ProjectUseCase) GetNamespaces(ctx context.Context) ([]string, error) {
	namespaces := make([]string, 0)
	kubeClientSet, err := GetKubeClientByInCluster()
	if err != nil {
		return nil, err
	}
	namespaceList, err := kubeClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.Name)
	}
	return namespaces, nil
}
