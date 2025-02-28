package repo

import (
	"context"
	"path/filepath"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	operatorApi "github.com/f-rambo/cloud-copilot/operator/api/v1alpha1"
	mateV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type cloudServiceClient struct {
	restClient rest.Interface
	ns         string
}

func NewClServiceClient(namespace string) (*cloudServiceClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
	}
	operatorApi.AddToScheme(scheme.Scheme)
	config.ContentConfig.GroupVersion = &operatorApi.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	restClient, err := rest.UnversionedRESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return &cloudServiceClient{
		restClient: restClient,
		ns:         namespace,
	}, nil
}

func (c *cloudServiceClient) List(ctx context.Context, opts mateV1.ListOptions) (*operatorApi.CloudServiceList, error) {
	result := operatorApi.CloudServiceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("cloudservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *cloudServiceClient) Get(ctx context.Context, name string, opts mateV1.GetOptions) (*operatorApi.CloudService, error) {
	result := operatorApi.CloudService{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Name(name).
		Resource("cloudservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *cloudServiceClient) Create(ctx context.Context, cs *operatorApi.CloudService) (*operatorApi.CloudService, error) {
	result := operatorApi.CloudService{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("cloudservices").
		Body(cs).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *cloudServiceClient) Update(ctx context.Context, cs *operatorApi.CloudService) (*operatorApi.CloudService, error) {
	result := operatorApi.CloudService{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Name(cs.Name).
		Resource("cloudservices").
		Body(cs).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *cloudServiceClient) Delete(ctx context.Context, name string, opts mateV1.DeleteOptions) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Name(name).
		Resource("cloudservices").
		Body(&opts).
		Do(ctx).
		Error()
}

func ConvertToCloudService(service *biz.Service, ci *biz.ContinuousIntegration, cd *biz.ContinuousDeployment) *operatorApi.CloudService {
	volumes := make([]operatorApi.Volume, 0)
	for _, v := range service.Volumes {
		volumes = append(volumes, operatorApi.Volume{
			Id:           v.Id,
			Name:         v.Name,
			Path:         v.MountPath,
			Storage:      v.Storage,
			StorageClass: v.StorageClass,
			ServiceId:    v.ServiceId,
		})
	}
	ports := make([]operatorApi.Port, 0)
	for _, v := range service.Ports {
		ports = append(ports, operatorApi.Port{
			Id:            v.Id,
			Name:          v.Name,
			ContainerPort: v.ContainerPort,
			Protocol:      v.Protocol,
			IngressPath:   v.IngressPath,
			ServiceId:     v.ServiceId,
		})
	}

	cs := &operatorApi.CloudService{
		TypeMeta: mateV1.TypeMeta{
			Kind:       "CloudService",
			APIVersion: "operator.cloud-copilot.com/v1alpha1",
		},
		ObjectMeta: mateV1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels:    utils.StringToMap(service.Lables),
		},
		Spec: operatorApi.CloudServiceSpec{
			Id:            service.Id,
			Image:         cd.Image,
			Replicas:      service.Replicas,
			RequestCPU:    service.RequestCpu,
			LimitCPU:      service.LimitCpu,
			RequestGPU:    service.RequestGpu,
			LimitGPU:      service.LimitGpu,
			RequestMemory: service.RequestMemory,
			LimitMemory:   service.LimitMemory,
			ConfigPath:    cd.ConfigPath,
			Config: map[string]string{
				filepath.Base(cd.ConfigPath): cd.Config,
			},
			Volumes: volumes,
			Ports:   ports,
		},
	}
	return cs
}

func SetServiceStatus(service *biz.Service, cs *operatorApi.CloudService) {
	service.Status = biz.ServiceStatus(cs.Status.Status)
}
