package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type nsCell corev1.Namespace

var Ns ns

type ns struct{}

type Nstotal struct {
	Name string
}

type NsFied struct {
	Name  string            `json:"name"`
	Label map[string]string `json:"label"`
}

type NsResp struct {
	Total int                `json:"total"`
	Items []corev1.Namespace `json:"items"`
}

func (n nsCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nsCell) GetName() string {
	return n.Name
}

func (n *ns) toCells(Ns []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(Ns))
	for i := range Ns {
		cells[i] = nsCell(Ns[i])
	}
	return cells
}

func (n *ns) fromCells(cells []DataCell) []corev1.Namespace {
	ns := make([]corev1.Namespace, len(cells))
	for i := range cells {
		ns[i] = corev1.Namespace(cells[i].(nsCell))
	}
	return ns
}

func (n *ns) GetNsNum() (t []string, err error) {
	var namespaceList []string
	NamespaceList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
	}
	for _, v := range NamespaceList.Items {
		namespaceList = append(namespaceList, v.Name)
	}

	return namespaceList, nil
}

func (n *ns) CreateNs(data *NsFied) (err error) {
	options := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Label,
		},
	}

	_, err = K8s.Clientset.CoreV1().Namespaces().Create(context.TODO(), options, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建namespace失败" + err.Error())
		return errors.New("创建namespace失败" + err.Error())
	}
	return nil
}

func (n *ns) DeleteNs(name string) (err error) {

	err = K8s.Clientset.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除namespace失败" + err.Error())
		return errors.New("删除namespace失败" + err.Error())
	}
	return nil
}

func (n *ns) GetNsDetails(name string) (namespace *corev1.Namespace, err error) {
	namespace, err = K8s.Clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取namespace详情失败" + err.Error())
		return nil, errors.New("获取namespace详情失败" + err.Error())
	}
	return namespace, nil
}

func (n *ns) UpdateNs(content string) (err error) {
	var namespace = &corev1.Namespace{}
	err = json.Unmarshal([]byte(content), namespace)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("namespace更新失败" + err.Error())
		return errors.New("namespace更新失败" + err.Error())
	}
	return nil
}
func (n *ns) GetNss(filterName string, limit, page int) (*NsResp, error) {

	nsList, err := K8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pod列表失败" + err.Error())
		return nil, errors.New("获取Pod列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: n.toCells(nsList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	pods := n.fromCells(data.GenericDataList)

	return &NsResp{
		Total: total,
		Items: pods,
	}, nil
}
