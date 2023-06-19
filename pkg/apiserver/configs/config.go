package configs

import (
	"github.com/myoperator/k8saggregatorapiserver/pkg/apis/myingress/v1beta1"
	"github.com/myoperator/k8saggregatorapiserver/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1beta1.AddToScheme(scheme))
}

func init() {

	Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	//metav1.AddToGroupVersion(Scheme, v1beta1.SchemeGroupVersion)

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// ExtraConfig holds custom controller configs
type ExtraConfig struct {
	// Place you custom configs here.
}

// Config defines the configs for the controller
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// TestServer contains state for a Kubernetes cluster master/api server.
type TestServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	// GenericConfig genericServer配置
	GenericConfig genericapiserver.CompletedConfig
	// ExtraConfig 自定义的额外配置
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}
	// FIXME 目前informer还无法使用
	//c.GenericConfig.SharedInformerFactory =
	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return CompletedConfig{&c}
}

// New returns a new instance of WardleServer from the given configs.
func (c completedConfig) New() (*TestServer, error) {

	// 生成apiserver底层的genericServer
	genericServer, err := c.GenericConfig.New("myapi", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	// 实例
	s := &TestServer{
		GenericAPIServer: genericServer,
	}

	// 重要步骤，创建和etcd交互
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(v1beta1.ApiGroup, Scheme, metav1.ParameterCodec, Codecs)
	v1beta1storage := map[string]rest.Storage{}
	// store.RESTInPeace store.NewREST 需要自己实现
	v1beta1storage["myingresses"] = store.RESTInPeace(store.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter))
	// 设置存储
	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1storage

	// 注册本apiserver路由
	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
