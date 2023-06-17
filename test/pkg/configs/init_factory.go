package configs

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

func InitClient() kubernetes.Interface {
	c, err := kubernetes.NewForConfig(K8sRestConfig())
	if err != nil {
		klog.Fatal("init k8s client error: ", err)
	}
	return c
}

var Factory informers.SharedInformerFactory

func InitInformer() {
	Factory = informers.NewSharedInformerFactoryWithOptions(InitClient(), 0)
	ingressInformer := Factory.Networking().V1().Ingresses()
	ingressInformer.Informer().AddEventHandler(IngressHandler{})
	stopC := make(chan struct{})
	Factory.Start(stopC)
	Factory.WaitForCacheSync(stopC)
}

var _ cache.ResourceEventHandler = &IngressHandler{}

type IngressHandler struct {
}

func (i IngressHandler) OnAdd(obj interface{}) {

}

func (i IngressHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (i IngressHandler) OnDelete(obj interface{}) {

}
