package service

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	config "k8sManager/config"
	"time"
)

var K8s k8s
var b = time.Now()
var Now = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", b.Year(), b.Month(), b.Day(), b.Hour(), b.Minute(), b.Second())

type k8s struct {
	Clientset *kubernetes.Clientset
}

func (k *k8s) InitClient() {

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config.ConfigPath)
	if err != nil {
		panic(err)
	}
	k.Clientset, err = kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err)
	}

}

// 父结构体
type dataSelector struct {
	GenericDataList   []DataCell
	dataSelectorQuery *DataSelectorQuery
}

// 排序结构体
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// 过滤分页父结构体
type DataSelectorQuery struct {
	FilterQuery     *FilterQuery
	PaginationQuery *PaginationQuery
}

// 过滤结构体，定义过滤关键字name
type FilterQuery struct {
	Name string
}

// 分页结构体，定义参数limit与page
type PaginationQuery struct {
	Limit int
	Page  int
}
