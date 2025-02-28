package biz

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/f-rambo/cloud-copilot/cluster-runtime/utils"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var ErrClusterNotFound error = errors.New("cluster not found")

type ClusterRuntimeConfigMapKey string

func (k ClusterRuntimeConfigMapKey) String() string {
	return string(k)
}

const (
	CloudCopilotNamespace = "cloud-copilot"

	ClusterInformation   ClusterRuntimeConfigMapKey = "cluster-info"
	NodegroupInformation ClusterRuntimeConfigMapKey = "nodegroup-info"
	NodeLableKey         ClusterRuntimeConfigMapKey = "node-lable"
)

func (c *Cluster) GetNodeGroup(nodeGroupId string) *NodeGroup {
	for _, nodeGroup := range c.NodeGroups {
		if nodeGroup.Id == nodeGroupId {
			return nodeGroup
		}
	}
	return nil
}

func (ng *NodeGroup) SetTargetSize(size int32) {
	ng.TargetSize = size
}

type ClusterUsecase struct {
	log *log.Helper
}

func NewClusterUseCase(logger log.Logger) *ClusterUsecase {
	c := &ClusterUsecase{
		log: log.NewHelper(logger),
	}
	return c
}

func GetKubeClientByKubeConfig(KubeConfigPaths ...string) (clientset *kubernetes.Clientset, err error) {
	var KubeConfigPath string
	if len(KubeConfigPaths) == 0 {
		KubeConfigPath = clientcmd.RecommendedHomeFile
	} else {
		KubeConfigPath = KubeConfigPaths[0]
	}
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes by kubeconfig client failed")
	}
	return client, nil
}

func GetKubeClientByRestConfig(masterIp, token, ca, key, cert string) (clientset *kubernetes.Clientset, err error) {
	if masterIp == "" || token == "" || ca == "" || key == "" || cert == "" {
		return nil, errors.New("invalid rest config")
	}
	config := &rest.Config{
		Host:        masterIp + ":6443",
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   []byte(ca),
			KeyData:  []byte(key),
			CertData: []byte(cert),
		},
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes by rest config client failed")
	}
	if client == nil {
		return nil, errors.New("get kubernetes by rest config client failed")
	}
	return client, nil
}

func GetDynamicClientByRestConfig(masterIp, token, ca, key, cert string) (*dynamic.DynamicClient, error) {
	config := &rest.Config{
		Host:        masterIp + ":6443",
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   []byte(ca),
			KeyData:  []byte(key),
			CertData: []byte(cert),
		},
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get dynamic client by rest config failed")
	}
	return client, nil
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

func (c *ClusterUsecase) CheckClusterInstalled(cluster *Cluster) error {
	_, err := GetKubeClientByRestConfig(cluster.ApiServerAddress, cluster.Token, cluster.CaData, cluster.KeyData, cluster.CertData)
	if err != nil {
		return ErrClusterNotFound
	}
	return nil
}

func (c *ClusterUsecase) CurrentCluster(ctx context.Context, cluster *Cluster) error {
	kubeClient, err := GetKubeClientByKubeConfig()
	if err != nil {
		return ErrClusterNotFound
	}
	err = c.getClusterInfo(ctx, kubeClient, cluster)
	if err != nil {
		return err
	}
	err = c.getNodes(ctx, kubeClient, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterUsecase) HandlerNodes(ctx context.Context, cluster *Cluster) error {
	clientset, err := GetKubeClientByRestConfig(cluster.ApiServerAddress, cluster.Token, cluster.CaData, cluster.KeyData, cluster.CertData)
	if err != nil {
		return err
	}
	for _, node := range cluster.Nodes {
		if node.Status != NodeStatus_NODE_DELETING {
			continue
		}

		pods, err := c.getPodsOnNode(ctx, clientset, node.Name)
		if err != nil {
			return fmt.Errorf("failed to get pods on node %s: %v", node.Name, err)
		}

		if err := c.evictPods(ctx, clientset, pods); err != nil {
			return fmt.Errorf("failed to evict pods: %v", err)
		}

		if err := c.waitForPodsToBeDeleted(ctx, clientset, pods, 5*time.Minute); err != nil {
			return fmt.Errorf("timeout waiting for pods to be deleted: %v", err)
		}

		err = clientset.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterUsecase) getClusterInfo(ctx context.Context, clientSet *kubernetes.Clientset, cluster *Cluster) error {
	configMap, err := clientSet.CoreV1().ConfigMaps(CloudCopilotNamespace).Get(ctx, ClusterInformation.String(), metav1.GetOptions{})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil
		}
		return err
	}
	clusterInfoString, ok := configMap.Data[ClusterInformation.String()]
	if !ok {
		return nil
	}
	err = utils.DeserializeFromBase64(clusterInfoString, cluster)
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
		for _, v := range node.Status.Addresses {
			if v.Address == "" {
				continue
			}
			if v.Type == corev1.NodeInternalIP {
				n.Ip = v.Address
			}
		}
		for _, v := range node.Status.Conditions {
			switch v.Type {
			case corev1.NodeReady:
				if v.Status == corev1.ConditionFalse {
					n.Status = NodeStatus_NODE_ERROR
					n.ErrorType = NodeErrorType_CLUSTER_ERROR
					n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", v.Reason, v.Message)
				}
			case corev1.NodeMemoryPressure:
				if v.Status == corev1.ConditionTrue {
					n.Status = NodeStatus_NODE_ERROR
					n.ErrorType = NodeErrorType_CLUSTER_ERROR
					n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", v.Reason, v.Message)
				}
			case corev1.NodeDiskPressure:
				if v.Status == corev1.ConditionTrue {
					n.Status = NodeStatus_NODE_ERROR
					n.ErrorType = NodeErrorType_CLUSTER_ERROR
					n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", v.Reason, v.Message)
				}
			case corev1.NodePIDPressure:
				if v.Status == corev1.ConditionTrue {
					n.Status = NodeStatus_NODE_ERROR
					n.ErrorType = NodeErrorType_CLUSTER_ERROR
					n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", v.Reason, v.Message)
				}
			case corev1.NodeNetworkUnavailable:
				if v.Status == corev1.ConditionTrue {
					n.Status = NodeStatus_NODE_ERROR
					n.ErrorType = NodeErrorType_CLUSTER_ERROR
					n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", v.Reason, v.Message)
				}
			default:
				n.Status = NodeStatus_NODE_ERROR
				n.ErrorType = NodeErrorType_CLUSTER_ERROR
				n.ErrorMessage += fmt.Sprintf("Reason: %s, Message: %s", corev1.ConditionUnknown, v.Message)
			}
		}
		for k, v := range node.Status.Capacity {
			gi, _ := toGiMi(&v)
			if k == corev1.ResourceCPU {
				n.NodeInfo += fmt.Sprintf("CPU(gi): %d ", gi)
			}
			if k == corev1.ResourceMemory {
				n.NodeInfo += fmt.Sprintf("Memory(gi): %d ", gi)
			}
			if k == corev1.ResourceEphemeralStorage {
				n.NodeInfo += fmt.Sprintf("Disk(gi): %d ", gi)
			}
			if strings.Contains(strings.ToUpper(k.String()), "GPU") {
				n.NodeInfo += fmt.Sprintf("GPU(%s): %d ", k.String(), gi)
			}
		}
		if clusterNodeIndex == -1 {
			cluster.Nodes = append(cluster.Nodes, n)
		} else {
			cluster.Nodes[clusterNodeIndex] = n
		}
	}
	return nil
}

func toGiMi(q *resource.Quantity) (int32, int32) {
	var gi, mi float64
	if q.Format == resource.BinarySI {
		gi = float64(q.Value()) / (1 << 30) // 1 Gi = 2^30 bytes
		mi = float64(q.Value()) / (1 << 20) // 1 Mi = 2^20 bytes
	}
	if q.Format == resource.DecimalSI {
		gi = float64(q.Value())
		mi = float64(q.Value()) * 1024
	}
	return int32(gi), int32(mi)
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
