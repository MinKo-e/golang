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

type nodeCell corev1.Node

var Node node

type node struct{}

type NodeResp struct {
	Total int           `json:"total"`
	Items []corev1.Node `json:"items"`
}

func (n nodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n nodeCell) GetName() string {
	return n.Name
}

func (n *node) toCells(Node []corev1.Node) []DataCell {
	cells := make([]DataCell, len(Node))
	for i := range Node {
		cells[i] = nodeCell(Node[i])
	}
	return cells
}

func (n *node) fromCells(cells []DataCell) []corev1.Node {
	ns := make([]corev1.Node, len(cells))
	for i := range cells {
		ns[i] = corev1.Node(cells[i].(nodeCell))
	}
	return ns
}

func (n *node) GetNodeDetails(name string) (nd *corev1.Node, err error) {
	nd, err = K8s.Clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取node详情失败" + err.Error())
		return nil, errors.New("获取node详情失败" + err.Error())
	}
	return nd, nil
}

func (n *node) UpdateNode(content string) (err error) {
	var nd = &corev1.Node{}
	err = json.Unmarshal([]byte(content), nd)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().Nodes().Update(context.TODO(), nd, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Node更新失败" + err.Error())
		return errors.New("Node更新失败" + err.Error())
	}
	return nil
}

func (n *node) GetNodes(filterName string, limit, page int) (*NodeResp, error) {

	ndList, err := K8s.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pod列表失败" + err.Error())
		return nil, errors.New("获取Pod列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: n.toCells(ndList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	node := n.fromCells(data.GenericDataList)

	return &NodeResp{
		Total: total,
		Items: node,
	}, nil
}

func (n *node) GetNodeRole() (t [][]string, err error) {
	var Master []string
	var Node []string

	NodeList, err := K8s.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for _, v1 := range NodeList.Items {
		node := true
		for k, _ := range v1.Labels {
			if k != "node-role.kubernetes.io/control-plane" || k != "node-role.kubernetes.io/master" {
				node = false
				continue
			}
		}
		if node {
			Node = append(Node, v1.Name)
		} else {
			Master = append(Master, v1.Name)
		}
	}
	return [][]string{Master, Node}, nil
}
