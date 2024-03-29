package configs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"
)

var K8sClient kubernetes.Interface

func init() {
	K8sClient = InitClient()
}

func K8sRestConfigInCluster() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal("init k8s rest configs error: ", err)
	}
	return config
}

func K8sRestConfig() *rest.Config {
	if os.Getenv("release") == "1" {
		klog.Infof("run in cluster")
		return K8sRestConfigInCluster()
	}

	klog.Infof("run outside cluster, debug mode")
	config, err := clientcmd.BuildConfigFromFlags("", "./resources/config")
	if err != nil {
		klog.Fatal("init k8s configs error: ", err)
	}

	return config
}
