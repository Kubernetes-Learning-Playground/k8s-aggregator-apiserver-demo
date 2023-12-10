package options

import (
	"github.com/spf13/pflag"
)

type ETCDOptions struct {
	// EtcdEndpoint 地址加端口
	EtcdEndpoint string
	// DefaultEtcdPathPrefix 存入etcd前缀
	// 如：/registry/xxx.xxx.com
	DefaultEtcdPathPrefix string
}

func NewETCDOptions() *ETCDOptions {
	return &ETCDOptions{}
}

func (o *ETCDOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.EtcdEndpoint, "etcd-server-endpoint", "http://127.0.0.1:2379",
		"ETCD service host address. Default to 127.0.0.1:2379.")
	fs.StringVar(&o.DefaultEtcdPathPrefix, "defaultEtcdPathPrefix", "/registry/myapi.jtthink.com",
		"Username for access to mysql service. Default to root.")
}

func (o *ETCDOptions) Complete() error {
	return nil
}

func (o *ETCDOptions) Validate() []error {
	var errs []error
	return errs
}
