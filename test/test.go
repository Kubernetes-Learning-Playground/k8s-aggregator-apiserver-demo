package main

import (
	"fmt"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apis/myingress/v1beta1"
	"github.com/myoperator/k8saggregatorapiserver/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	netutils "k8s.io/utils/net"
	"log"
	"net"
)

const (
	RemoteKubeConfig      = "./resources/config"
	defaultEtcdPathPrefix = "/registry/myapi.jtthink.com"
)

var (
	// Scheme 定义了资源序列化和反序列化的方法以及资源类型和版本的对应关系
	Scheme = runtime.NewScheme()
	//  编解码器工厂
	Codecs = serializer.NewCodecFactory(Scheme)
)

// TODO: 修改注解与改成入参配置
// 推荐配置 模板
func getRcOpt() *genericoptions.RecommendedOptions {
	rc := genericoptions.NewRecommendedOptions(
		defaultEtcdPathPrefix,
		Codecs.LegacyCodec(v1beta1.SchemeGroupVersion), //JSON格式的编码器
	)

	rc.SecureServing.BindPort = 6443
	rc.SecureServing.ServerCert = genericoptions.GeneratableKeyCert{
		CertDirectory: "./cert",
		PairName:      "aaserver",
	}

	rc.Etcd.StorageConfig.Transport.ServerList = []string{"http://127.0.0.1:2379"}

	rc.CoreAPI.CoreAPIKubeconfigPath = RemoteKubeConfig
	rc.Authentication.RemoteKubeConfigFile = RemoteKubeConfig
	rc.Authorization.RemoteKubeConfigFile = RemoteKubeConfig

	rc.CoreAPI = nil
	rc.Authentication = nil
	rc.Authorization = nil
	rc.Admission = nil

	err := rc.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"0.0.0.0",
		nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")})

	if err != nil {
		log.Fatalln(err)
	}
	return rc
}

//func main() {
//
//	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Group: "", Version: "v1"})
//	// TODO: keep the generic API server from wanting this
//	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
//	//这些东西和常规的POD、Deployment不同，仅仅是表示资源的字符串标识
//	Scheme.AddUnversionedTypes(unversioned,
//		&metav1.Status{},
//		&metav1.APIVersions{},
//		&metav1.APIGroupList{},
//		&metav1.APIGroup{},
//		&metav1.APIResourceList{},
//	)
//
//	//从这里开始看 ，把我们的myingress加入
//	err := v1beta1.AddToScheme(Scheme)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	//生成 controller 的推荐配置（默认配置)
//	config := genericapiserver.NewRecommendedConfig(Codecs)
//
//	err = getRcOpt().ApplyTo(config) //模板赋值
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	//定义存储 (自定义 存储)
//	//resources := map[string]rest.Storage{
//	//	"myingresses": store.NewMyStore(v1beta1.SchemeGroupResource, true,
//	//		rest.NewDefaultTableConvertor(v1beta1.SchemeGroupResource)),
//	//}
//
//	gvi := v1beta1.SchemeGroupVersion
//	gvi.Version = runtime.APIVersionInternal
//	Scheme.AddKnownTypes(gvi, &v1beta1.MyIngress{}, &v1beta1.MyIngressList{})
//	agi := genericapiserver.NewDefaultAPIGroupInfo(
//		v1beta1.SchemeGroupVersion.Group,
//		Scheme,
//		metav1.ParameterCodec, Codecs)
//	completeConfig := config.Complete()
//	resources := map[string]rest.Storage{
//		"myingresses": store.RESTInPeace(store.NewREST(Scheme,
//			completeConfig.RESTOptionsGetter)),
//	}
//
//	//设置存储
//	agi.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = resources
//
//
//	config.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(v1beta1.GetOpenAPIDefinitions,
//		openapi.NewDefinitionNamer(Scheme))
//	//configs.OpenAPIConfig.Info.Title = "MyIngress"
//	//configs.OpenAPIConfig.Info.Version = "v1"
//
//
//
//	// controller
//	server, err := completeConfig.
//		New("myapi", genericapiserver.NewEmptyDelegate())
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	err = server.InstallAPIGroup(&agi)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	stopCh := genericapiserver.SetupSignalHandler()
//	err = server.PrepareRun().Run(stopCh)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//}

func main() {

	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	//这些东西和常规的POD、Deployment不同，仅仅是表示资源的字符串标识
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)

	//从这里开始看 ，把我们的myingress加入
	err := v1beta1.AddToScheme(Scheme)
	if err != nil {
		log.Fatalln(err)
	}
	gvi := v1beta1.SchemeGroupVersion
	gvi.Version = runtime.APIVersionInternal
	Scheme.AddKnownTypes(gvi, &v1beta1.MyIngress{}, &v1beta1.MyIngressList{})

	agi := genericapiserver.NewDefaultAPIGroupInfo(
		v1beta1.SchemeGroupVersion.Group,
		Scheme,
		metav1.ParameterCodec, Codecs)

	//生成 controller 的推荐配置（默认配置)
	config := genericapiserver.NewRecommendedConfig(Codecs)

	err = getRcOpt().ApplyTo(config) //模板赋值
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("aaa")
	config.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(v1beta1.GetOpenAPIDefinitions,
		openapi.NewDefinitionNamer(Scheme))
	//config.OpenAPIConfig.Info.Title = "MyIngress"
	//config.OpenAPIConfig.Info.Version = "v1"
	completeConfig := config.Complete()
	//定义存储 (自定义 存储)
	//resources := map[string]rest.Storage{
	//	"myingresses": store.NewMyStore(v1beta1.SchemeGroupResource, true,
	//		rest.NewDefaultTableConvertor(v1beta1.SchemeGroupResource)),
	//}

	resources := map[string]rest.Storage{
		"myingresses": store.RESTInPeace(store.NewREST(Scheme,
			completeConfig.RESTOptionsGetter)),
	}
	//设置存储
	agi.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = resources
	// controller
	fmt.Println("ddd")
	server, err := completeConfig.
		New("myapi", genericapiserver.NewEmptyDelegate())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("dddsss")
	err = server.InstallAPIGroup(&agi)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("runrun")
	stopCh := genericapiserver.SetupSignalHandler()
	err = server.PrepareRun().Run(stopCh)
	if err != nil {
		log.Fatalln(err)
	}

}
