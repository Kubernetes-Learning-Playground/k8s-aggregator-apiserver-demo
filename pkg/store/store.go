package store

import (
	"github.com/practice/k8s_aggregator_apiserver/pkg/apis/myingress/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// 代替etcd，使用内存方式
var IngressMap map[string][]*v1beta1.MyIngress

func init() {
	IngressMap = make(map[string][]*v1beta1.MyIngress)

	test := &v1beta1.MyIngress{}
	test.Name = "test-ingress"
	test.Namespace = "default"
	test.Spec.Path = "test-path"

	createMap(test)

	test2 := &v1beta1.MyIngress{}
	test2.Name = "测试Ingress2"
	test2.Namespace = "default"
	test2.Spec.Path = "/jtthink/abc"
	test2.Spec.Host = "abc.jtthink.com"
	createMap(test2)

}

func createMap(ingress *v1beta1.MyIngress) {
	ingress.CreationTimestamp = metav1.NewTime(time.Now())
	if _, ok := IngressMap[ingress.Namespace]; !ok {
		IngressMap[ingress.Namespace] = []*v1beta1.MyIngress{}
	}
	IngressMap[ingress.Namespace] = append(IngressMap[ingress.Namespace], ingress)
}

func findIngressByNamespace(ns string) []*v1beta1.MyIngress {
	if list, ok := IngressMap[ns]; !ok {
		IngressMap[ns] = []*v1beta1.MyIngress{}
		return IngressMap[ns]
	} else {
		return list
	}

}

func ListIngressMap(ns string) *v1beta1.MyIngressList {
	list := v1beta1.NewMyIngressList()
	list.Items = findIngressByNamespace(ns)
	return list
}
