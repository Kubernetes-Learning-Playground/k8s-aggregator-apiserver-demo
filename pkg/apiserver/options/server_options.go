package options

import (
	"github.com/spf13/pflag"
)

type ServerOptions struct {
	// RemoteKubeConfig kube-config文件
	RemoteKubeConfig string
	// CertFiles 签发证书目录
	CertFiles string
	// PairName 签发证书文件名，如签发证书为
	// aaserver.crt aaserver.key 则填入 aaserver
	PairName string
	// Port aggregator-apiserver端口
	Port int
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{}
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.RemoteKubeConfig, "kube-config", "./resources/config",
		"")
	fs.StringVar(&o.CertFiles, "certFiles", "./cert",
		"")
	fs.StringVar(&o.PairName, "pairName", "aaserver",
		"")

	fs.IntVar(&o.Port, "server-port", 6443,
		"")

}

func (o *ServerOptions) Complete() error {
	return nil
}

func (o *ServerOptions) Validate() []error {
	var errs []error

	return errs
}
