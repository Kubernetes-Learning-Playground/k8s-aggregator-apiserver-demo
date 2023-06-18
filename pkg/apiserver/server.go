package apiserver

import (
	"fmt"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apis/myingress/v1beta1"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apiserver/configs"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/informers"
	"net"
)

// TestServerOptions contains state for master/api server
type TestServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	SharedInformerFactory informers.SharedInformerFactory
	StdOut                io.Writer
	StdErr                io.Writer
}

// NewTestServerOptions returns a new WardleServerOptions
func NewTestServerOptions(out, errOut io.Writer, opt *configs.ParameterConfig) *TestServerOptions {
	o := &TestServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			opt.DefaultEtcdPathPrefix,
			// json格式解码器
			configs.Codecs.LegacyCodec(v1beta1.SchemeGroupVersion),
		),

		StdOut: out,
		StdErr: errOut,
	}
	o.RecommendedOptions.SecureServing.BindPort = opt.Port
	o.RecommendedOptions.SecureServing.ServerCert = genericoptions.GeneratableKeyCert{
		CertDirectory: opt.CertFiles,
		PairName:      opt.PairName,
	}
	o.RecommendedOptions.CoreAPI.CoreAPIKubeconfigPath = opt.RemoteKubeConfig
	o.RecommendedOptions.Authentication.RemoteKubeConfigFile = opt.RemoteKubeConfig
	o.RecommendedOptions.Authorization.RemoteKubeConfigFile = opt.RemoteKubeConfig
	o.RecommendedOptions.Etcd.StorageConfig.Transport.ServerList = []string{opt.EtcdEndpoint}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(v1beta1.SchemeGroupVersion, schema.GroupKind{Group: v1beta1.ApiGroup})
	return o
}

// NewCommandStartTestServer provides a CLI handler for 'start master' command
// with a default TestServerOptions.
func NewCommandStartTestServer(defaults *TestServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a test API server",
		Long:  "Launch a test API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunTestServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)

	return cmd
}

// Validate validates TestServerOptions
func (o TestServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *TestServerOptions) Complete() error {
	// register admission plugins
	//banflunder.Register(o.RecommendedOptions.Admission.Plugins)

	// add admission plugins to the RecommendedPluginOrder
	//o.RecommendedOptions.Admission.RecommendedPluginOrder = append(o.RecommendedOptions.Admission.RecommendedPluginOrder, "BanFlunder")

	return nil
}

// Config returns configs for the api server given TestServerOptions
func (o *TestServerOptions) Config() (*configs.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	// FIXME: 使用code-generator 只有生成v1beta版本资源对象，无法使用完整的informer
	//o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)
	//client, err := clientset.NewForConfig(c.LoopbackClientConfig)
	//if err != nil {
	//	return nil, err
	//}
	//informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
	//o.SharedInformerFactory = informerFactory
	//o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
	//	client, err := clientset.NewForConfig(c.LoopbackClientConfig)
	//	if err != nil {
	//		return nil, err
	//	}
	//	informerFactory := informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
	//	o.SharedInformerFactory = informerFactory
	//	return []admission.PluginInitializer{wardleinitializer.New(informerFactory)}, nil
	//}

	serverConfig := genericapiserver.NewRecommendedConfig(configs.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(v1beta1.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(configs.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = v1beta1.ResourceKind
	serverConfig.OpenAPIConfig.Info.Version = v1beta1.ApiVersion

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &configs.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   configs.ExtraConfig{},
	}
	return config, nil
}

// RunTestServer starts a new WardleServer given TestServerOptions
func (o TestServerOptions) RunTestServer(stopCh <-chan struct{}) error {
	cc, err := o.Config()
	if err != nil {
		return err
	}

	server, err := cc.Complete().New()
	if err != nil {
		return err
	}

	//server.GenericAPIServer.AddPostStartHookOrDie("start-sample-server-informers", func(context genericapiserver.PostStartHookContext) error {
	//	cc.GenericConfig.SharedInformerFactory.Start(context.StopCh)
	//	o.SharedInformerFactory.Start(context.StopCh)
	//	return nil
	//})

	// TODO: 心跳检测健康机制
	//go func() {
	//	h := &healthz.Handler{
	//		Checks: map[string]healthz.Checker{
	//			"healthz": healthz.Ping,
	//		},
	//	}
	//	klog.Info(fmt.Sprintf("0.0.0.0:%d", options.HealthPort))
	//	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", options.HealthPort), h); err != nil {
	//		klog.Fatalf("Failed to start healthz endpoint: %v", err)
	//	}
	//}()

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
