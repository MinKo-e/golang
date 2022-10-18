package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

var NetworkIngress networkIngress

type networkIngress struct{}

type NetworkIngresssResp struct {
	Total int                 `json:"total"`
	Items []networkv1.Ingress `json:"items"`
}

// 定义networkIngress结构体参数，用于接收前端数据
type NetworkIngressFied struct {
	Name        string                 `form:"name" binding:"required" json:"name"`
	Namespace   string                 `form:"namespace" binding:"required" json:"namespace"`
	Labels      map[string]string      `form:"labels" json:"labels"`
	Annotations map[string]string      `form:"annotations" json:"annotations"`
	TLS         bool                   `form:"tls" json:"tls"`
	SAName      string                 `json:"sa_name"`
	Hosts       map[string][]*HttpPath `json:"hosts"`
}

// {"host":{}}
type HttpPath struct {
	Path        string `json:"path" binding:"required"`
	PathType    string `json:"path_type"`
	ServiceName string `json:"service_name" binding:"required"`
	ServicePort int    `json:"service_port" binding:"required"`
}

// {"test1":{ {path},{path} }}

type NetworkIngresstotal struct {
	NetworkIngressNum int
	Namespace         string
}

// 定义corev1.networkIngress数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type networkIngressCell networkv1.Ingress

func (n networkIngressCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n networkIngressCell) GetName() string {
	return n.Name
}

func (n *networkIngress) toCells(networkIngresss []networkv1.Ingress) []DataCell {
	cells := make([]DataCell, len(networkIngresss))
	for i := range networkIngresss {
		cells[i] = networkIngressCell(networkIngresss[i])
	}
	return cells
}

func (n *networkIngress) fromCells(cells []DataCell) []networkv1.Ingress {
	networkIngresss := make([]networkv1.Ingress, len(cells))
	for i := range cells {
		networkIngresss[i] = networkv1.Ingress(cells[i].(networkIngressCell))
	}
	return networkIngresss
}

func (n *networkIngress) GetNetworkIngressNum() (t []NetworkIngresstotal, err error) {
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
		networkIngress_list, err := K8s.Clientset.NetworkingV1().Ingresses(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取NetworkIngress列表失败" + err.Error())
			return nil, err

		}
		t = append(t, NetworkIngresstotal{
			NetworkIngressNum: len(networkIngress_list.Items),
			Namespace:         v,
		})
	}

	return t, nil
}

func (n *networkIngress) GetNetworkIngresss(filterName, namespace string, limit, page int) (*NetworkIngresssResp, error) {

	networkIngressList, err := K8s.Clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取NetworkIngress列表失败" + err.Error())
		return nil, errors.New("获取NetworkIngress列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: n.toCells(networkIngressList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	networkIngresss := n.fromCells(data.GenericDataList)

	return &NetworkIngresssResp{
		Total: total,
		Items: networkIngresss,
	}, nil
}

func (n *networkIngress) GetNetworkIngressDetails(name, namespace string) (networkIngress *networkv1.Ingress, err error) {
	networkIngress, err = K8s.Clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取NetworkIngress详情失败" + err.Error())
		return nil, errors.New("获取NetworkIngress详情失败" + err.Error())
	}
	return networkIngress, nil
}

func (n *networkIngress) DeleteNetworkIngress(name, namespace string) (err error) {
	err = K8s.Clientset.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除NetworkIngress失败" + err.Error())
		return errors.New("删除NetworkIngress失败" + err.Error())
	}
	return nil
}

func (p *networkIngress) CreateNetworkIngress(networkIngressstruct *NetworkIngressFied) (err error) {

	var ingressHttpPath []networkv1.HTTPIngressPath
	var ingressRules []networkv1.IngressRule
	var ingressdefaultpathtype = networkv1.PathTypePrefix
	var ingressDefaultClass = "nginx"
	options := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:              networkIngressstruct.Name,
			Namespace:         networkIngressstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            networkIngressstruct.Labels,
			Annotations:       networkIngressstruct.Annotations,
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: &ingressDefaultClass,
			TLS:              nil,
		},
	}

	for key, value := range networkIngressstruct.Hosts {
		ir := networkv1.IngressRule{
			Host: key,
			IngressRuleValue: networkv1.IngressRuleValue{HTTP: &networkv1.HTTPIngressRuleValue{
				Paths: nil,
			}}}
		for _, HttpPath := range value {

			hip := networkv1.HTTPIngressPath{
				Path:     HttpPath.Path,
				PathType: &ingressdefaultpathtype,
				Backend: networkv1.IngressBackend{
					Service: &networkv1.IngressServiceBackend{
						Name: HttpPath.ServiceName,
						Port: networkv1.ServiceBackendPort{
							Number: int32(HttpPath.ServicePort),
						},
					},
				},
			}
			ingressHttpPath = append(ingressHttpPath, hip)
		}
		ir.IngressRuleValue.HTTP.Paths = ingressHttpPath
		ingressRules = append(ingressRules, ir)
	}
	options.Spec.Rules = ingressRules
	fmt.Printf("#%v\n", ingressRules)
	_, err = K8s.Clientset.NetworkingV1().Ingresses(networkIngressstruct.Namespace).Create(context.TODO(), options, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建NetworkIngress失败" + err.Error())
		return errors.New("创建NetworkIngress失败" + err.Error())
	}
	return nil
}

func (p *networkIngress) UpdateNetworkIngress(namespace, content string) (err error) {
	var networkIngress = &networkv1.Ingress{}
	err = json.Unmarshal([]byte(content), networkIngress)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), networkIngress, metav1.UpdateOptions{})

	if err != nil {
		logrus.Error("NetworkIngress更新失败" + err.Error())
		return errors.New("NetworkIngress更新失败" + err.Error())
	}
	return nil
}
