package main

import (
	"flag"
	"os"

	"k8s.io/klog/v2"

	"github.com/myoperator/k8saggregatorapiserver/pkg/apiserver"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
)

func main() {
	// 日志相关
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := apiserver.NewTestServerOptions(os.Stdout, os.Stderr)
	cmd := apiserver.NewCommandStartTestServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
