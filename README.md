### 基于k8s提供的aggregator-apiserver实现自定义资源
#### 项目思路：
使用k8s提供的aggregator-apiserver，实现集群内自定义资源对象(不同于crd+controller形式扩展)。


### 部署方式：
注：部署前需要先创建aggregator api-server的CA文件，否则无法挂载。详细可以参考网上内容。
1. 编译可执行文件 or 打镜像(注意：这里需要自行适配，deploy.yaml中适用镜像方式)
```bigquery
# 打镜像
docker build -t aggregator-apiserver:v1 .
# 编译可执行文件
docker run --rm -it -v /root/k8s-aggregator-controller:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./myapi .
```
2. 部署应用
- 需要部署APIService资源对象
- deployment aggregator-apiserver服务
- service 对外访问服务
```bigquery
[root@VM-0-16-centos k8s-aggregator-apiserver]# cd yaml/
[root@VM-0-16-centos yaml]# kubectl apply -f .
apiservice.apiregistration.k8s.io/v1beta1.apis.jtthink.com unchanged
deployment.apps/myapi unchanged
service/myapi unchanged
```
查看服务对象
```bigquery
[root@VM-0-16-centos ~]# kubectl get pods | grep myapi
myapi-c74bfb6f5-gcbht                       2/2     Running             0                  24h
```

查看自定义资源
```bigquery
[root@VM-0-16-centos ~]# kubectl api-resources | grep myingress
myingresses                       mi           apis.jtthink.com/v1beta1               true         MyIngress
```
先查看自定义资源
```bigquery
[root@VM-0-16-centos ~]# kubectl get myingresses.apis.jtthink.com
No resources found in default namespace.
```
操作自定义资源
```bigquery
[root@VM-0-16-centos yaml]# ls
api.yaml  deploy.yaml  etcd.yaml  myingress.yaml  mypod.yaml  rbac.yaml
[root@VM-0-16-centos yaml]# kubectl apply -f myingress.yaml
myingress.apis.jtthink.com/test-myingress created
[root@VM-0-16-centos yaml]# kubectl get myingresses.apis.jtthink.com
NAME             CREATED AT
test-myingress   2023-06-19T15:41:25Z
[root@VM-0-16-centos yaml]# kubectl delete -f myingress.yaml
myingress.apis.jtthink.com "test-myingress" deleted
[root@VM-0-16-centos yaml]# kubectl get myingresses.apis.jtthink.com
No resources found in default namespace.
```