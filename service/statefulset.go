package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type stsCell appsv1.StatefulSet

var Sts sts

type sts struct{}

type Ststotal struct {
	DsNum     int
	Namespace string
}

type StsFied struct {
	Name     string            `json:"name"`
	Label    map[string]string `json:"label"`
	Replicas int32             `json:"replicas"`
	PodFied
}

type StsResp struct {
	Total int                  `json:"total"`
	Items []appsv1.StatefulSet `json:"items"`
}

func (s stsCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s stsCell) GetName() string {
	return s.Name
}

func (s *sts) toCells(sts []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(sts))
	for i := range sts {
		cells[i] = stsCell(sts[i])
	}
	return cells
}

func (s *sts) fromCells(cells []DataCell) []appsv1.StatefulSet {
	ds := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		ds[i] = appsv1.StatefulSet(cells[i].(stsCell))
	}
	return ds
}

func (s *sts) GetStsNum() (t []Ststotal, err error) {
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
		sts_list, err := K8s.Clientset.AppsV1().StatefulSets(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取StatefulSet列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Ststotal{
			DsNum:     len(sts_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (s *sts) ScaleSts(name, namespace string, replicas int) (replica int32, err error) {

	scale, err := K8s.Clientset.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取副本数失败" + err.Error())
		return 0, errors.New("获取副本数失败" + err.Error())
	}
	scale.Spec.Replicas = int32(replicas)
	newscale, err := K8s.Clientset.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("更新副本数失败" + err.Error())
		return 0, errors.New("更新副本数失败" + err.Error())
	}
	return newscale.Spec.Replicas, nil
}

func (s *sts) GetSts(filterName, namespace string, limit, page int) (*StsResp, error) {

	stsList, err := K8s.Clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取StatefulSet列表失败" + err.Error())
		return nil, errors.New("获取StatefulSet列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: s.toCells(stsList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	sts := s.fromCells(data.GenericDataList)

	return &StsResp{
		Total: total,
		Items: sts,
	}, nil
}

func (s *sts) GetStsDetails(name, namespace string) (sts *appsv1.StatefulSet, err error) {
	sts, err = K8s.Clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取StatefulSet详情失败" + err.Error())
		return nil, errors.New("获取StatefulSet详情失败" + err.Error())
	}
	return sts, nil
}

func (s *sts) DeleteSts(name, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Statefulset失败" + err.Error())
		return errors.New("删除Statefulset失败" + err.Error())
	}
	return nil
}

func (s *sts) UpdateSts(namespace, content string) (err error) {
	var sts = &appsv1.StatefulSet{}
	err = json.Unmarshal([]byte(content), sts)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), sts, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("damonset更新失败" + err.Error())
		return errors.New("damonset更新失败" + err.Error())
	}
	return nil
}

func (s *sts) CreateSts(data *StsFied) (err error) {

	container := corev1.PodSpec{
		Containers: []corev1.Container{
			{Name: data.CName, Image: data.CImage},
		},
		NodeSelector:       data.NodeSelector,
		ServiceAccountName: data.ServiceAccountName,
		NodeName:           data.NodeName,
	}
	if data.IName != "" && data.IImage != "" {
		container.InitContainers = []corev1.Container{
			{Name: data.IName, Image: data.IImage},
		}
	}
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              data.PName,
			Namespace:         data.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            data.Labels,
		},
		Spec: container,
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
			Replicas: &data.Replicas,
		},
	}
	_, err = K8s.Clientset.AppsV1().StatefulSets(data.Namespace).Create(context.TODO(), sts, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建StatefulSet失败" + err.Error())
		return errors.New("创建StatefulSet失败" + err.Error())
	}
	return nil
}

func (s *sts) RestartSts(name, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{"name": name,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": Now,
							}},
						},
					},
				},
			},
		},
	}
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logrus.Error("序列化数据失败" + err.Error())
		return errors.New("序列化数据失败" + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().StatefulSets(namespace).Patch(context.TODO(), name,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logrus.Error("重启StatefulSet失败" + err.Error())
		return errors.New("重启StatefulSet失败" + err.Error())
	}
	return nil
}
