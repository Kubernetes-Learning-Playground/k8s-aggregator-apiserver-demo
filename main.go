package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/practice/k8s_aggregator_apiserver/pkg/apis/myingress/v1beta1"
	"github.com/practice/k8s_aggregator_apiserver/pkg/builders"
	"github.com/practice/k8s_aggregator_apiserver/pkg/configs"
	"github.com/practice/k8s_aggregator_apiserver/pkg/store"
	"github.com/practice/k8s_aggregator_apiserver/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

//把 xx=xx,xx=xxx  解析为一个map
func parseLabelQuery(query string) map[string]string {
	m := make(map[string]string)
	if query == "" {
		return m
	}
	qs := strings.Split(query, ",")
	if len(qs) == 0 {
		return m
	}
	for _, q := range qs {
		qPair := strings.Split(q, "=")
		if len(qPair) == 2 {
			m[qPair[0]] = qPair[1]
		}
	}
	return m
}

var rootJson = `
{
  "kind":"APIResourceList",
  "apiVersion":"v1",
  "groupVersion":"apis.jtthink.com/v1beta1",
  "resources":[
     {"name":"mypods","singularName":"mypod","shortNames":["mp"],"namespaced":true,"kind":"MyPod","verbs":["get","list"]}
  ]}
`
var podsListv2 = `
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items":[
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
`
var podsListv1 = `
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items":[
    {
	  "metadata": {
        "name": "testpod1-v1",
        "namespace": "default"
       }
    },
    {
	  "metadata": {
        "name": "testpod2-v1",
        "namespace": "default"
       }
    }
   ]
}
`
var podDetail = `
{
  "kind": "MyPod",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {"name":"{name}","namespace":"{namespace}"},
  "spec":{"属性":"你懂的"},
  "columnDefinitions": [
        {
            "name": "Name",
            "type": "string"
        },
        {
            "name": "Created At",
            "type": "date"
        }
    ]
}
`

var (
	ROOTURL = fmt.Sprintf("/apis/%s/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version)
	//根据NS 获取 myingress列表
	ListByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
)

func main() {

	configs.InitInformer()

	r := gin.New()

	r.Use(func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
		fmt.Println(c.Request.URL.RawQuery)
		c.Next()
	})

	// 根 必须实现
	r.GET(ROOTURL, func(c *gin.Context) {
		c.JSON(200, builders.ApiResourceList())
	})

	//r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/myingresses", func(c *gin.Context) {
	//	ns := c.Param("ns")
	//	res := utils.ConvertToTable(store.ListIngressMap(ns))
	//	c.JSON(200, res)
	//})

	r.GET(ListByNS_URL, func(c *gin.Context) {
		ns := c.Param("ns")
		rr, _ := store.NewClientStore().ListByNamespaceOrAll(ns)
		res := utils.ConvertToTable(rr)
		c.JSON(200, res)
	})

	// 模拟用的，不需要了。
	//r.GET("/apis/apis.jtthink.com/v1beta1", func(c *gin.Context) {
	//	c.Header("content-type","application/json")
	//	c.String(200, rootJson)
	//})

	//列表: 根据namespace来查询pods
	r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/mypods", func(c *gin.Context) {
		c.Header("content-type", "application/json")

		// 解析出query 参数(labelQuery)
		labelQueryMap := parseLabelQuery(c.Query("labelSelector"))
		json := ""
		if v, ok := labelQueryMap["version"]; ok {
			if v == "1" {
				json = strings.Replace(podsListv1, "default", c.Param("ns"), -1)
			}
		}

		if json == "" {
			json = strings.Replace(podsListv2, "default", c.Param("ns"), -1)
		}

		c.String(200, json)
	})

	//列表 list
	r.GET("/apis/apis.jtthink.com/v1beta1/mypods", func(c *gin.Context) {
		c.Header("content-type", "application/json")
		json := strings.Replace(podsListv1, "default", "all", -1)
		c.String(200, json)
	})

	//详细 （根据ns)
	r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/mypods/:name", func(c *gin.Context) {
		//c.Header("content-type","application/json")
		//json:=strings.Replace(podDetail,"{namespace}",c.Param("ns"),-1)
		//json=strings.Replace(json,"{name}",c.Param("name"),-1)
		//
		//
		//c.String(200,json)
		t := metav1.Table{}
		t.Kind = "Table"
		t.APIVersion = "meta.k8s.io/v1"
		t.ColumnDefinitions = []metav1.TableColumnDefinition{
			{Name: "name", Type: "string"},
			{Name: "命名空间", Type: "string"},
			{Name: "状态", Type: "string"},
		}
		t.Rows = []metav1.TableRow{
			{Cells: []interface{}{c.Param("name"),
				c.Param("ns"), "准备好了"},
			},
		}
		c.JSON(200, t)

	})

	//if err := r.RunTLS(":8443",
	//	"/etc/kubernetes/pki/aaserver.crt", "/etc/kubernetes/pki/aaserver.key"); err != nil {
	//	log.Fatalln(err)
	//}
	if err := r.RunTLS(":8443", "./cert/aaserver.crt", "./cert/aaserver.key"); err != nil {
		log.Fatalln(err)
	}

}
