### kubernetes aggregator api-server 练习

### 部署方式：
注：部署前需要先创建aggregator api-server的CA文件，否则无法挂载。详细可以参考网上内容。
1. 编译可执行文件
```bigquery
docker run --rm -it -v /root/k8s-aggregator-apiserver:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./myapi .
```
2. 部署应用
```bigquery
[root@VM-0-16-centos k8s-aggregator-apiserver]# cd yaml/
[root@VM-0-16-centos yaml]# kubectl apply -f .
apiservice.apiregistration.k8s.io/v1beta1.apis.jtthink.com unchanged
deployment.apps/myapi unchanged
service/myapi unchanged
```
```bigquery
[root@VM-0-16-centos yaml]# kubectl get pods | grep myapi
myapi-c5d8674b8-6p62v                  1/1     Running   0              26m
```

3. demo版本，没有正式的业务逻辑。项目处于待更新的阶段
目前可以请求server路由，如下。
```bigquery
[root@VM-0-16-centos yaml]# kubectl get mypods
NAME          AGE
testpod1-v2   <unknown>
testpod2-v2   <unknown>
[root@VM-0-16-centos yaml]# kubectl get mypods testpod1-v2
NAME          命名空间      状态
testpod1-v2   default   准备好了
[root@VM-0-16-centos yaml]# kubectl get mypods testpod1-v2 -oyaml
apiVersion: meta.k8s.io/v1
columnDefinitions:
- description: ""
  format: ""
  name: name
  priority: 0
  type: string
- description: ""
  format: ""
  name: 命名空间
  priority: 0
  type: string
- description: ""
  format: ""
  name: 状态
  priority: 0
  type: string
kind: Table
metadata: {}
rows:
- cells:
  - testpod1-v2
  - default
  - 准备好了
  object: null
```   
```bigquery
[root@VM-0-16-centos yaml]# kubectl get --raw  /apis/apis.jtthink.com/v1beta1   | jq
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "apis.jtthink.com/v1beta1",
  "resources": [
    {
      "name": "mypods",
      "singularName": "mypod",
      "shortNames": [
        "mp"
      ],
      "namespaced": true,
      "kind": "MyPod",
      "verbs": [
        "get",
        "list"
      ]
    }
  ]
}
[root@VM-0-16-centos yaml]# kubectl get --raw  /apis/apis.jtthink.com/v1beta1/namespaces/default/mypods   | jq
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items": [
    {
      "metadata": {
        "name": "testpod1-v2",
        "namespace": "default"
      }
    },
    {
      "metadata": {
        "name": "testpod2-v2",
        "namespace": "default"
      }
    }
  ]
}
```
