package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var ErrClusterNotFound error = errors.New("cluster not found")

type ClusterRuntimeConfigMapKey string

func (k ClusterRuntimeConfigMapKey) String() string {
	return string(k)
}

const (
	ClusterInformation   ClusterRuntimeConfigMapKey = "cluster-info"
	NodegroupInformation ClusterRuntimeConfigMapKey = "nodegroup-info"
	NodeLableKey         ClusterRuntimeConfigMapKey = "node-lable"
)

type ClusterUsecase struct {
	log *log.Helper
}

func NewClusterUseCase(logger log.Logger) *ClusterUsecase {
	c := &ClusterUsecase{
		log: log.NewHelper(logger),
	}
	return c
}

func (c *ClusterUsecase) CreateYAMLFile(ctx context.Context, dynamicClient *dynamic.DynamicClient, namespace, resource, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}
	defer file.Close()
	decoder := yaml.NewYAMLOrJSONDecoder(file, 1024)
	for {
		unstructuredObj := &unstructured.Unstructured{}
		if err := decoder.Decode(unstructuredObj); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return errors.Wrap(err, "decode yaml failed")
		}
		gvr := unstructuredObj.GroupVersionKind().GroupVersion().WithResource(resource)
		resourceClient := dynamicClient.Resource(gvr).Namespace(namespace)
		_, err = resourceClient.Create(ctx, unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to create resource")
		}
	}
	return nil
}

func (c *ClusterUsecase) CurrentCluster(ctx context.Context, cluster *Cluster) (*Cluster, error) {
	kubeClient, err := GetKubeClientByInCluster()
	if err != nil {
		return cluster, ErrClusterNotFound
	}
	versionInfo, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		return cluster, err
	}
	cluster.Version = versionInfo.String()
	err = c.getClusterInfo(ctx, kubeClient, cluster)
	if err != nil {
		return cluster, err
	}
	err = c.getNodes(ctx, kubeClient, cluster)
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (c *ClusterUsecase) HandlerNodes(ctx context.Context, cluster *Cluster) (*Cluster, error) {
	clientset, err := GetKubeClientByRestConfig(cluster.MasterIp, cluster.Token, cluster.CaData, cluster.KeyData, cluster.CertData)
	if err != nil {
		return cluster, err
	}
	for _, node := range cluster.Nodes {
		if node.Status != NodeStatus_NODE_DELETING {
			continue
		}

		pods, err := c.getPodsOnNode(ctx, clientset, node.Name)
		if err != nil {
			return cluster, fmt.Errorf("failed to get pods on node %s: %v", node.Name, err)
		}

		if err := c.evictPods(ctx, clientset, pods); err != nil {
			return cluster, fmt.Errorf("failed to evict pods: %v", err)
		}

		if err := c.waitForPodsToBeDeleted(ctx, clientset, pods, 5*time.Minute); err != nil {
			return cluster, fmt.Errorf("timeout waiting for pods to be deleted: %v", err)
		}

		err = clientset.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
		if err != nil {
			return cluster, err
		}
	}
	return cluster, nil
}

func (c *ClusterUsecase) MigrateToCluster(ctx context.Context, cluster *Cluster) (*Cluster, error) {
	clientset, err := GetKubeClientByRestConfig(cluster.MasterIp, cluster.Token, cluster.CaData, cluster.KeyData, cluster.CertData)
	if err != nil {
		return cluster, err
	}
	serverName := "cloud-copilot"
	labels := map[string]string{"app.kubernetes.io/cluster.name": cluster.Name, "app.kubernetes.io/name": serverName}
	namspaceObj, err := clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serverName,
			Labels: labels,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return cluster, err
	}
	serviceAccountObj, err := clientset.CoreV1().ServiceAccounts(namspaceObj.Namespace).Create(ctx, &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-app-service-account", serverName),
			Namespace: namspaceObj.Namespace,
			Labels:    labels,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return cluster, err
	}
	_, err = clientset.RbacV1().ClusterRoleBindings().Create(ctx, &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-app-cluster-admin-binding", serverName),
			Labels: labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountObj.Name,
				Namespace: serviceAccountObj.Namespace,
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return cluster, err
	}
	var replicas int32 = 1
	cloudCopilotPodSpec := corev1.PodSpec{
		ServiceAccountName: serviceAccountObj.Name,
		HostNetwork:        true,
		Containers: []corev1.Container{
			{
				Name:  fmt.Sprintf("%s-app-container", serverName),
				Image: "frambo/cloud-copilot:v0.0.1",
				Env:   []corev1.EnvVar{},
				Ports: []corev1.ContainerPort{
					{ContainerPort: 8000},
					{ContainerPort: 9000},
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
		},
	}
	deploymentObj, err := clientset.AppsV1().Deployments(namspaceObj.Namespace).Create(ctx, &appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-app-deployment", serverName),
			Namespace: namspaceObj.Namespace,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: cloudCopilotPodSpec,
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return cluster, err
	}
	c.log.Info(deploymentObj.Name)
	_, err = clientset.CoreV1().Services(namspaceObj.Namespace).Create(ctx, &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-app-service", serverName),
			Namespace: namspaceObj.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: fmt.Sprintf("%s.%s.com", serverName, cluster.Name),
			Ports: []corev1.ServicePort{
				{
					Port:       8000,
					TargetPort: intstr.FromInt(8000),
				},
				{
					Port:       9000,
					TargetPort: intstr.FromInt(9000),
				},
			},
			Selector: labels,
		},
		Status: corev1.ServiceStatus{},
	}, metav1.CreateOptions{})
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (c *ClusterUsecase) getClusterInfo(ctx context.Context, clientSet *kubernetes.Clientset, cluster *Cluster) error {
	configMap, err := clientSet.CoreV1().ConfigMaps("kube-system").Get(ctx, ClusterInformation.String(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	if _, ok := configMap.Data[ClusterInformation.String()]; !ok {
		return nil
	}
	err = json.Unmarshal([]byte(configMap.Data[ClusterInformation.String()]), cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterUsecase) getNodes(ctx context.Context, clientSet *kubernetes.Clientset, cluster *Cluster) error {
	nodeRes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodeRes.Items {
		n := &Node{}
		clusterNodeIndex := -1
		for index, v := range cluster.Nodes {
			if v.Name == node.Name {
				n = v
				clusterNodeIndex = index
				break
			}
		}
		n.Name = node.Name
		for _, v := range node.Status.Addresses {
			if v.Address == "" {
				continue
			}
			if v.Type == "InternalIP" {
				n.InternalIp = v.Address
			}
			if v.Type == "ExternalIP" {
				n.ExternalIp = v.Address
			}
		}
		n.Status = NodeStatus_NODE_UNSPECIFIED
		for _, v := range node.Status.Conditions {
			if v.Status == corev1.ConditionStatus(corev1.NodeReady) {
				n.Status = NodeStatus_NODE_RUNNING
			}
		}
		nodeLables, err := json.Marshal(node)
		if err != nil {
			return err
		}
		err = json.Unmarshal(nodeLables, &n)
		if err != nil {
			return err
		}
		n.Labels = string(nodeLables)
		if clusterNodeIndex == -1 {
			cluster.Nodes = append(cluster.Nodes, n)
		} else {
			cluster.Nodes[clusterNodeIndex] = n
		}
	}
	return nil
}

func (c *ClusterUsecase) getPodsOnNode(ctx context.Context, clientset *kubernetes.Clientset, nodeName string) ([]corev1.Pod, error) {
	podList, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return nil, err
	}

	var podsToEvict []corev1.Pod
	for _, pod := range podList.Items {
		if c.isMirrorPod(&pod) || c.isDaemonSetPod(&pod) {
			continue
		}
		podsToEvict = append(podsToEvict, pod)
	}
	return podsToEvict, nil
}

func (c *ClusterUsecase) isMirrorPod(pod *corev1.Pod) bool {
	_, exists := pod.Annotations[corev1.MirrorPodAnnotationKey]
	return exists
}

func (c *ClusterUsecase) isDaemonSetPod(pod *corev1.Pod) bool {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "DaemonSet" {
			return true
		}
	}
	return false
}

func (c *ClusterUsecase) evictPods(ctx context.Context, clientset *kubernetes.Clientset, pods []corev1.Pod) error {
	for _, pod := range pods {
		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
		}
		err := clientset.PolicyV1().Evictions(eviction.Namespace).Evict(ctx, eviction)
		if err != nil {
			return errors.Wrap(err, "failed to evict pod")
		}
	}
	return nil
}

func (c *ClusterUsecase) waitForPodsToBeDeleted(ctx context.Context, clientset *kubernetes.Clientset, pods []corev1.Pod, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return wait.PollUntilContextTimeout(ctx, time.Second, timeout, true, func(ctx context.Context) (done bool, err error) {
		for _, pod := range pods {
			_, err := clientset.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err == nil {
				return false, nil
			}
		}
		return true, nil
	})
}
