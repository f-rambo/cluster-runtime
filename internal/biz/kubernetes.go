package biz

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func GetKubeClientByInCluster() (clientset *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes by in cluster client failed")
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
