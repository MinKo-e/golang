package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sManager/service"
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

func (n *node) toCells(Node []corev1.Node) []service.DataCell {
	cells := make([]service.DataCell, len(Node))
	for i := range Node {
		cells[i] = nodeCell(Node[i])
	}
	return cells
}

func (n *node) fromCells(cells []service.DataCell) []corev1.Node {
	ns := make([]corev1.Node, len(cells))
	for i := range cells {
		ns[i] = corev1.Node(cells[i].(nodeCell))
	}
	return ns
}

func (n *node) GetNodeDetails(name string) (nd *corev1.Node, err error) {
	nd, err = service.K8s.Clientset.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
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
	_, err = service.K8s.Clientset.CoreV1().Nodes().Update(context.TODO(), nd, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Node更新失败" + err.Error())
		return errors.New("Node更新失败" + err.Error())
	}
	return nil
}
