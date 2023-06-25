package configs

// ParameterConfig 命令行入参配置
type ParameterConfig struct {
	// RemoteKubeConfig kube-config文件
	RemoteKubeConfig string
	// EtcdEndpoint 地址加端口
	EtcdEndpoint string
	// DefaultEtcdPathPrefix 存入etcd前缀
	// 如：/registry/xxx.xxx.com
	DefaultEtcdPathPrefix string
	// CertFiles 签发证书目录
	CertFiles string
	// PairName 签发证书文件名，如签发证书为
	// aaserver.crt aaserver.key 则填入 aaserver
	PairName string
	// Port aggregator-apiserver端口
	Port int
}
