package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

var NetworkSvc networkSvc

type networkSvc struct{}

type NetworkSvcsResp struct {
	Total int              `json:"total"`
	Items []corev1.Service `json:"items"`
}

// 定义networkSvc结构体参数，用于接收前端数据
type NetworkSvcFied struct {
	Name       string            `form:"name" binding:"required" json:"name"`
	Namespace  string            `form:"namespace" binding:"required" json:"namespace"`
	Labels     map[string]string `form:"labels" json:"labels"`
	Type       string            `json:"type" form:"type" binding:"required"`
	Selector   map[string]string `form:"selector" json:"selector"`
	NodePort   int               `form:"node_port" json:"node_port"`
	TargetPort int               `form:"target_port" json:"target_port"`
	Port       int               `form:"port" json:"port"`
}

type NetworkSvctotal struct {
	NetworkSvcNum int
	Namespace     string
}

// 定义corev1.networkSvc数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type networkSvcCell corev1.Service

func (n networkSvcCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n networkSvcCell) GetName() string {
	return n.Name
}

func (n *networkSvc) toCells(networkSvcs []corev1.Service) []DataCell {
	cells := make([]DataCell, len(networkSvcs))
	for i := range networkSvcs {
		cells[i] = networkSvcCell(networkSvcs[i])
	}
	return cells
}

func (n *networkSvc) fromCells(cells []DataCell) []corev1.Service {
	networkSvcs := make([]corev1.Service, len(cells))
	for i := range cells {
		networkSvcs[i] = corev1.Service(cells[i].(networkSvcCell))
	}
	return networkSvcs
}

func (n *networkSvc) GetNetworkSvcNum() (t []NetworkSvctotal, err error) {
	var namespaceList []string
	NamespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
	}
	for _, v := range NamespaceList.Items {
		namespaceList = append(namespaceList, v.Name)
	}
	fmt.Println(namespaceList)

	for _, v := range namespaceList {
		networkSvc_list, err := K8s.Clientset.CoreV1().Services(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取NetworkSvc列表失败" + err.Error())
			return nil, err

		}
		t = append(t, NetworkSvctotal{
			NetworkSvcNum: len(networkSvc_list.Items),
			Namespace:     v,
		})
	}

	return t, nil
}

func (n *networkSvc) GetNetworkSvcs(filterName, namespace string, limit, page int) (*NetworkSvcsResp, error) {

	networkSvcList, err := K8s.Clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取NetworkSvc列表失败" + err.Error())
		return nil, errors.New("获取NetworkSvc列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: n.toCells(networkSvcList.Items), dataSelectorQuery: &DataSelectorQuery{
		FilterQuery: &FilterQuery{Name: filterName},
		PaginationQuery: &PaginationQuery{
			Limit: limit,
			Page:  page,
		},
	}}

	//先过滤，后排序分页
	filterQuery := selectoerQuery.Filter()
	total := len(filterQuery.GenericDataList)
	data := filterQuery.Sort().Paging()
	networkSvcs := n.fromCells(data.GenericDataList)

	return &NetworkSvcsResp{
		Total: total,
		Items: networkSvcs,
	}, nil
}

func (n *networkSvc) GetNetworkSvcDetails(name, namespace string) (networkSvc *corev1.Service, err error) {
	networkSvc, err = K8s.Clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取NetworkSvc详情失败" + err.Error())
		return nil, errors.New("获取NetworkSvc详情失败" + err.Error())
	}
	return networkSvc, nil
}

func (n *networkSvc) DeleteNetworkSvc(name, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除NetworkSvc失败" + err.Error())
		return errors.New("删除NetworkSvc失败" + err.Error())
	}
	return nil
}

func (p *networkSvc) CreateNetworkSvc(networkSvcstruct *NetworkSvcFied) (err error) {
	options := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:              networkSvcstruct.Name,
			Namespace:         networkSvcstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            networkSvcstruct.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     int32(networkSvcstruct.Port),
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: int32(networkSvcstruct.TargetPort),
					},
				},
			},
			Selector: networkSvcstruct.Selector,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	if networkSvcstruct.Type != "" {
		if networkSvcstruct.Type == "node_port" {
			options.Spec.Type = corev1.ServiceTypeNodePort
			options.Spec.Ports[0].NodePort = int32(networkSvcstruct.NodePort)
		} else if networkSvcstruct.Type == "headless" {
			options.Spec.ClusterIP = ""
		} else if networkSvcstruct.Type == "loadbalancer" {
			options.Spec.Type = corev1.ServiceTypeLoadBalancer
		}

	}
	_, err = K8s.Clientset.CoreV1().Services(networkSvcstruct.Namespace).Create(context.TODO(), options, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建NetworkSvc失败" + err.Error())
		return errors.New("创建NetworkSvc失败" + err.Error())
	}
	return nil
}

func (p *networkSvc) UpdateNetworkSvc(namespace, content string) (err error) {
	var networkSvc = &corev1.Service{}
	err = json.Unmarshal([]byte(content), networkSvc)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().Services(namespace).Update(context.TODO(), networkSvc, metav1.UpdateOptions{})

	if err != nil {
		logrus.Error("NetworkSvc更新失败" + err.Error())
		return errors.New("NetworkSvc更新失败" + err.Error())
	}
	return nil
}
