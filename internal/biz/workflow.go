package biz

import (
	"context"
	"encoding/json"
	"fmt"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfcommon "github.com/argoproj/argo-workflows/v3/workflow/common"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type WorkflowV1Alpha1Interface interface {
	Workflows(namespace string) WorkflowInterface
}

type WorkflowV1Alpha1Client struct {
	restClient rest.Interface
}

/*
	argo workflow 资源
	customresourcedefinition.apiextensions.k8s.io/clusterworkflowtemplates.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/cronworkflows.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workflowartifactgctasks.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workfloweventbindings.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workflows.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workflowtaskresults.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workflowtasksets.argoproj.io created
	customresourcedefinition.apiextensions.k8s.io/workflowtemplates.argoproj.io created
	serviceaccount/argo created
	serviceaccount/argo-server created
	role.rbac.authorization.k8s.io/argo-role created
	clusterrole.rbac.authorization.k8s.io/argo-aggregate-to-admin created
	clusterrole.rbac.authorization.k8s.io/argo-aggregate-to-edit created
	clusterrole.rbac.authorization.k8s.io/argo-aggregate-to-view created
	clusterrole.rbac.authorization.k8s.io/argo-cluster-role created
	clusterrole.rbac.authorization.k8s.io/argo-server-cluster-role created
	rolebinding.rbac.authorization.k8s.io/argo-binding created
	clusterrolebinding.rbac.authorization.k8s.io/argo-binding created
	clusterrolebinding.rbac.authorization.k8s.io/argo-server-binding created
	configmap/workflow-controller-configmap created
	service/argo-server created
	priorityclass.scheduling.k8s.io/workflow-controller created
	deployment.apps/argo-server created
	deployment.apps/workflow-controller created
*/

/*
	argo cli
	kubectl create -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/main/examples/hello-world.yaml
	kubectl get wf -n argo
	kubectl get wf hello-world-xxx -n argo
	kubectl get po -n argo --selector=workflows.argoproj.io/workflow=hello-world-xxx
	kubectl logs hello-world-yyy -c main -n argo
*/

var (
	resourceName = "workflows"
)

type WorkflowInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*wfv1.WorkflowList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*wfv1.Workflow, error)
	Create(context.Context, *wfv1.Workflow) (*wfv1.Workflow, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string) error
	// ...
}

type workflowClient struct {
	restClient rest.Interface
	ns         string
}

func (c *workflowClient) List(ctx context.Context, opts metav1.ListOptions) (*wfv1.WorkflowList, error) {
	result := wfv1.WorkflowList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *workflowClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*wfv1.Workflow, error) {
	result := wfv1.Workflow{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resourceName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *workflowClient) Create(ctx context.Context, wf *wfv1.Workflow) (*wfv1.Workflow, error) {
	result := wfv1.Workflow{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(resourceName).
		Body(wf).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *workflowClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(resourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *workflowClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(resourceName).
		Name(name).
		Do(ctx).
		Error()
}

// unmarshalWorkflows unmarshals the input bytes as either json or yaml
func UnmarshalWorkflow(wfStr string, strict bool) (wfv1.Workflow, error) {
	wfs, err := UnmarshalWorkflows(wfStr, strict)
	if err != nil {
		return wfv1.Workflow{}, err
	}
	for _, v := range wfs {
		return v, nil
	}
	return wfv1.Workflow{}, nil
}

func UnmarshalWorkflows(wfStr string, strict bool) ([]wfv1.Workflow, error) {
	wfBytes := []byte(wfStr)
	return wfcommon.SplitWorkflowYAMLFile(wfBytes, strict)
}

func getKubeConfig() (config *rest.Config, err error) {
	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return
		}
	}
	return
}

func newForConfig(c *rest.Config) (*WorkflowV1Alpha1Client, error) {
	wfv1.AddToScheme(scheme.Scheme)
	config := *c
	config.ContentConfig.GroupVersion = &wfv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &WorkflowV1Alpha1Client{restClient: client}, nil
}

func (c *WorkflowV1Alpha1Client) Workflows(namespace string) WorkflowInterface {
	return &workflowClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func GetApiVersion() string {
	return fmt.Sprintf("%s/%s", wfv1.WorkflowSchemaGroupVersionKind.Group, wfv1.WorkflowSchemaGroupVersionKind.Version)
}

func GetKind() string {
	return wfv1.WorkflowSchemaGroupVersionKind.Kind
}

func GetDefaultWorklfows(_ context.Context, business, technology, args string) ([]byte, error) {
	wf := golangCiDefaultWorklfows()
	return json.Marshal(wf)
}

func golangCiDefaultWorklfows() *wfv1.Workflow {
	typeMeta := metav1.TypeMeta{
		APIVersion: GetApiVersion(),
		Kind:       GetKind(),
	}
	objMeta := metav1.ObjectMeta{
		GenerateName: "golang-ci-default-cloud-copilot-",
	}
	argu := wfv1.Arguments{
		Parameters: []wfv1.Parameter{
			{Name: "name"},
			{Name: "repo"},
			{Name: "registry"},
			{Name: "registry_user"},
			{Name: "registry_pwd"},
			{Name: "version"},
			{Name: "branch"},
			{Name: "tag"},
		},
	}
	voTem := apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "work",
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.VolumeResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
	volumeMounts := []apiv1.VolumeMount{
		{
			Name:      "work",
			MountPath: "/app",
		},
	}
	// 重点内容
	templates := []wfv1.Template{
		{
			Name: "main",
			Steps: []wfv1.ParallelSteps{
				{
					Steps: []wfv1.WorkflowStep{
						{
							Name:     "default-clone",
							Template: "clone-dag",
						},
						{
							Name:     "default-dep",
							Template: "deps-dag",
						},
						{
							Name:     "default-build",
							Template: "build-dag",
						},
						{
							Name:     "default-image",
							Template: "image-dag",
						},
					},
				},
			},
		},
		{
			Name: "clone-dag",
			DAG: &wfv1.DAGTemplate{
				Tasks: []wfv1.DAGTask{
					{
						Name:     "default-clone-dag-task",
						Template: "clone",
					},
				},
			},
		},
		{
			Name: "deps-dag",
			DAG: &wfv1.DAGTemplate{
				Tasks: []wfv1.DAGTask{
					{
						Name:     "default-deps-dag-task",
						Template: "deps",
						// Dependencies: []string{"clone"},
					},
				},
			},
		},
		{
			Name: "build-dag",
			DAG: &wfv1.DAGTemplate{
				Tasks: []wfv1.DAGTask{
					{
						Name:     "default-build-dag-task",
						Template: "build",
						// Dependencies: []string{"deps"},
					},
				},
			},
		},
		{
			Name: "image-dag",
			DAG: &wfv1.DAGTemplate{
				Tasks: []wfv1.DAGTask{
					{
						Name:     "default-image-dag-task",
						Template: "image",
						// Dependencies: []string{"deps"},
					},
				},
			},
		},
		{
			Name: "image",
			Container: &apiv1.Container{
				VolumeMounts: volumeMounts,
				Name:         "image",
				Image:        "moby/buildkit:latest",
				Command: []string{
					"sh",
					"-c",
				},
				Args: []string{
					`
buildctl --creds {{workflow.parameters.registry_user}}:{{workflow.parameters.registry_pwd}} \
--registry {{workflow.parameters.registry}} \
build \
--frontend dockerfile.v0 \
--local context=. \
--local dockerfile=. \
--output type=image,name={{workflow.parameters.name}}:{{workflow.parameters.version}},push=true`,
				},
			},
		},
		{
			Name: "clone",
			Container: &apiv1.Container{
				VolumeMounts: volumeMounts,
				Name:         "clone",
				Image:        "bitnami/git:latest",
				Command: []string{
					"sh",
					"-c",
				},
				Args: []string{
					`
if [ -z "{{workflow.parameters.tag}}" ]
then
git clone -v -b "{{workflow.parameters.branch}}" --single-branch --depth 1 "{{workflow.parameters.repo}}" .
else
git clone -v -b "{{workflow.parameters.tag}}" --single-branch --depth 1 "{{workflow.parameters.repo}}" .
fi`,
				},
			},
		},
		{
			Name: "deps",
			Container: &apiv1.Container{
				VolumeMounts: volumeMounts,
				Name:         "deps",
				Image:        "golang:1.23.3",
				Command: []string{
					"sh",
					"-c",
				},
				Args: []string{
					`make generate`,
				},
			},
		},
		{
			Name: "build",
			Container: &apiv1.Container{
				VolumeMounts: volumeMounts,
				Name:         "build",
				Image:        "golang:1.23.3",
				Command: []string{
					"sh",
					"-c",
				},
				Args: []string{
					`make build`,
				},
			},
		},
	}
	spec := wfv1.WorkflowSpec{
		Arguments:  argu,
		Entrypoint: "main",
		// OnExit:               "image",
		VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{voTem},
		Templates:            templates,
	}
	wf := &wfv1.Workflow{
		TypeMeta:   typeMeta,
		ObjectMeta: objMeta,
		Spec:       spec,
	}
	return wf
}
