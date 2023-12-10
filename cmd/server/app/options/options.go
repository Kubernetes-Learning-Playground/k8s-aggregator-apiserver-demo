package options

import (
	"flag"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apiserver/options"
	"github.com/pkg/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
)

type Options struct {
	Server *options.ServerOptions
	ETCD   *options.ETCDOptions
	Logs   *logs.Options
}

func NewOptions() *Options {
	return &Options{
		Server: options.NewServerOptions(),
		ETCD:   options.NewETCDOptions(),
		Logs:   logs.NewOptions(),
	}
}

func (o *Options) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	fss.FlagSet("generic").AddGoFlagSet(flag.CommandLine)

	logs.AddGoFlags(flag.CommandLine)

	o.Server.AddFlags(fss.FlagSet("server"))
	o.ETCD.AddFlags(fss.FlagSet("etcd"))
	return fss
}

func (o *Options) Complete() error {

	if err := o.Server.Complete(); err != nil {
		return err
	}
	if err := o.ETCD.Complete(); err != nil {
		return err
	}

	return nil
}

func (o *Options) Validate() error {
	var errs []error

	errs = append(errs, o.Server.Validate()...)
	errs = append(errs, o.ETCD.Validate()...)

	if len(errs) == 0 {
		return nil
	}

	wrapped := errors.New("options validate error")
	for _, err := range errs {
		wrapped = errors.WithMessage(wrapped, err.Error())
	}
	return wrapped
}
