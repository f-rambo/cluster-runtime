package repo

import (
	"context"
	"strings"

	argoworkflow "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfcommon "github.com/argoproj/argo-workflows/v3/workflow/common"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/internal/biz"
	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

const (
	ArgoWorkflowEntryTmpName = "main"

	ArgoWorkflowServiceAccount = "argo-server"
)

type workflowClient struct {
	restClient rest.Interface
	ns         string
}

func NewArgoWorkflowClient(namespace string) (*workflowClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, err
		}
	}
	wfv1.AddToScheme(scheme.Scheme)
	config.ContentConfig.GroupVersion = &wfv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}
	return &workflowClient{restClient: client, ns: namespace}, nil
}

func (c *workflowClient) List(ctx context.Context, opts metav1.ListOptions) (*wfv1.WorkflowList, error) {
	result := wfv1.WorkflowList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(argoworkflow.WorkflowPlural).
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
		Resource(argoworkflow.WorkflowPlural).
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
		Resource(argoworkflow.WorkflowPlural).
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
		Resource(argoworkflow.WorkflowPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

func (c *workflowClient) Delete(ctx context.Context, name string) error {
	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(argoworkflow.WorkflowPlural).
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

func ConvertToArgoWorkflow(w *biz.Workflow) *wfv1.Workflow {
	var deleteWfSecond int32 = 600
	argoWf := &wfv1.Workflow{
		TypeMeta: metav1.TypeMeta{
			APIVersion: argoworkflow.APIVersion,
			Kind:       argoworkflow.WorkflowKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      w.Name,
			Namespace: w.Namespace,
			Labels:    utils.StringToMap(w.Lables),
		},
		Spec: wfv1.WorkflowSpec{
			ServiceAccountName: ArgoWorkflowServiceAccount,
			Entrypoint:         ArgoWorkflowEntryTmpName,
			Templates:          make([]wfv1.Template, 0),
			Volumes: []apiv1.Volume{
				{
					Name: w.GetWorkdirName(),
					VolumeSource: apiv1.VolumeSource{
						PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
							ClaimName: w.GetStorageName(),
						},
					},
				},
			},
			TTLStrategy: &wfv1.TTLStrategy{
				SecondsAfterCompletion: &deleteWfSecond,
				SecondsAfterSuccess:    &deleteWfSecond,
				SecondsAfterFailure:    &deleteWfSecond,
			},
		},
	}
	for _, step := range w.WorkflowSteps {
		for _, task := range step.WorkflowTasks {
			task.Name = strings.ToLower(task.Name)
			argoWf.Spec.Templates = append(argoWf.Spec.Templates, wfv1.Template{
				Name: task.Name,
				Container: &apiv1.Container{
					Name:       task.Name,
					Image:      step.Image,
					Command:    []string{"sh", "-c"},
					Args:       []string{task.TaskCommand},
					WorkingDir: w.GetWorkdir(),
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      w.GetWorkdirName(),
							MountPath: w.GetWorkdir(),
						},
					},
				},
			})
		}
	}
	DAGTasks := make([]wfv1.DAGTask, 0)
	firstStep := w.GetFirstStep()
	if firstStep == nil {
		return nil
	}
	firstTask := firstStep.GetFirstTask()
	if firstTask == nil {
		return nil
	}

	DAGTasks = append(DAGTasks, wfv1.DAGTask{
		Name:     firstTask.Name,
		Template: firstTask.Name,
	})
	currentStep := firstStep
	currentTask := firstTask
	prevTaskName := currentTask.Name

	for {
		nextTask := currentStep.GetNextTask(currentTask)

		if nextTask == nil {
			nextStep := w.GetNextStep(currentStep)
			if nextStep == nil {
				break
			}
			currentStep = nextStep
			nextTask = currentStep.GetFirstTask()
			if nextTask == nil {
				break
			}
		}

		DAGTasks = append(DAGTasks, wfv1.DAGTask{
			Name:         nextTask.Name,
			Template:     nextTask.Name,
			Dependencies: []string{prevTaskName},
		})

		prevTaskName = nextTask.Name
		currentTask = nextTask
	}
	argoWf.Spec.Templates = append(argoWf.Spec.Templates, wfv1.Template{
		Name: ArgoWorkflowEntryTmpName,
		DAG:  &wfv1.DAGTemplate{Tasks: DAGTasks},
	})
	return argoWf
}

func SetWorkflowStatus(argoWf *wfv1.Workflow, wf *biz.Workflow) {
	for nodeName, node := range argoWf.Status.Nodes {
		task := wf.GetTask(nodeName)
		if task == nil {
			continue
		}
		switch node.Phase {
		case wfv1.NodePending, wfv1.NodeRunning:
			task.Status = biz.WorkfloStatus_Pending
		case wfv1.NodeSucceeded:
			task.Status = biz.WorkfloStatus_Success
		case wfv1.NodeFailed, wfv1.NodeError:
			task.Status = biz.WorkfloStatus_Failure
		default:
		}
	}
}
