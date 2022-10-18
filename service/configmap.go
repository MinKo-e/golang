package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

var Configmap configmap

type configmap struct{}

type ConfigmapsResp struct {
	Total int                `json:"total"`
	Items []corev1.ConfigMap `json:"items"`
}

// 定义configmap结构体参数，用于接收前端数据
type ConfigmapFied struct {
	Name      string            `form:"name" binding:"required" json:"name"`
	Namespace string            `form:"namespace" binding:"required" json:"namespace"`
	Label     map[string]string `form:"label" json:"labels"`
	Data      map[string]string `form:"data" json:"data"`
}

type Configmaptotal struct {
	ConfigmapNum int
	Namespace    string
}

// 定义corev1.configmap数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type configmapCell corev1.ConfigMap

func (c configmapCell) GetCreation() time.Time {
	return c.CreationTimestamp.Time
}

func (c configmapCell) GetName() string {
	return c.Name
}

func (c *configmap) toCells(configmaps []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(configmaps))
	for i := range configmaps {
		cells[i] = configmapCell(configmaps[i])
	}
	return cells
}

func (c *configmap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configmaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configmaps[i] = corev1.ConfigMap(cells[i].(configmapCell))
	}
	return configmaps
}

func (c *configmap) GetConfigmapNum() (t []Configmaptotal, err error) {
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
		configmap_list, err := K8s.Clientset.CoreV1().ConfigMaps(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取Configmap列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Configmaptotal{
			ConfigmapNum: len(configmap_list.Items),
			Namespace:    v,
		})
	}

	return t, nil
}

func (c *configmap) GetConfigmaps(filterName, namespace string, limit, page int) (*ConfigmapsResp, error) {

	configmapList, err := K8s.Clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Configmap列表失败" + err.Error())
		return nil, errors.New("获取Configmap列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: c.toCells(configmapList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	configmaps := c.fromCells(data.GenericDataList)

	return &ConfigmapsResp{
		Total: total,
		Items: configmaps,
	}, nil
}

func (c *configmap) GetConfigmapDetails(name, namespace string) (configmap *corev1.ConfigMap, err error) {
	configmap, err = K8s.Clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Configmap详情失败" + err.Error())
		return nil, errors.New("获取Configmap详情失败" + err.Error())
	}
	return configmap, nil
}

func (c *configmap) DeleteConfigmap(name, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Configmap失败" + err.Error())
		return errors.New("删除Configmap失败" + err.Error())
	}
	return nil
}

func (c *configmap) CreateConfigmap(configmapstruct *ConfigmapFied) (err error) {
	option := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:              configmapstruct.Name,
			Namespace:         configmapstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            configmapstruct.Label,
		},
		Data: configmapstruct.Data,
	}
	_, err = K8s.Clientset.CoreV1().ConfigMaps(configmapstruct.Namespace).Create(context.TODO(), option, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建Configmap失败" + err.Error())
		return errors.New("创建Configmap失败" + err.Error())
	}
	return nil
}

func (c *configmap) UpdateConfigmap(namespace, content string) (err error) {
	var configmap = &corev1.ConfigMap{}
	err = json.Unmarshal([]byte(content), configmap)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configmap, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Configmap更新失败" + err.Error())
		return errors.New("Configmap更新失败" + err.Error())
	}
	return nil
}

func (c *configmap) GetConfigmapData(name, namespace string) (map[string]string, error) {

	data, err := c.GetConfigmapDetails(name, namespace)
	if err != nil {
		logrus.Error("获取Configmap详情失败" + err.Error())
		return nil, errors.New("获取Configmap详情失败" + err.Error())
	}

	return data.Data, nil
}
