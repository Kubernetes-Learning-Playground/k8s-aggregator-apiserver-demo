package main

import (
	"flag"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apiserver"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apiserver/configs"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	"os"
)

var (
	etcdEndpoint          string
	remoteKubeConfig      string
	defaultEtcdPathPrefix string
	certFiles             string
	pairName              string
	port                  int
)

func main() {

	flag.StringVar(&etcdEndpoint, "etcd-server-endpoint", "http://127.0.0.1:2379", "")
	flag.StringVar(&remoteKubeConfig, "kube-config", "./resources/config", "")
	flag.StringVar(&defaultEtcdPathPrefix, "defaultEtcdPathPrefix", "/registry/myapi.jtthink.com", "")
	flag.StringVar(&certFiles, "certFiles", "./cert", "")
	flag.StringVar(&pairName, "pairName", "aaserver", "")
	flag.IntVar(&port, "server-port", 8443, "")
	flag.Parse()

	// 配置文件
	opt := &configs.ParameterConfig{
		EtcdEndpoint:          etcdEndpoint,
		RemoteKubeConfig:      remoteKubeConfig,
		DefaultEtcdPathPrefix: defaultEtcdPathPrefix,
		CertFiles:             certFiles,
		PairName:              pairName,
		Port:                  port,
	}

	// 日志相关
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := apiserver.NewTestServerOptions(os.Stdout, os.Stderr, opt)
	cmd := apiserver.NewCommandStartTestServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
