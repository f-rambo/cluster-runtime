package repo

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesClientSet() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return kubernetes.NewForConfig(inClusterConfig)
	}
	return kubernetes.NewForConfig(config)
}
