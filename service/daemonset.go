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

type dsCell appsv1.DaemonSet

var Ds ds

type ds struct{}

type Dstotal struct {
	DsNum     int
	Namespace string
}

type DsFied struct {
	Name  string            `json:"name"`
	Label map[string]string `json:"label"`
	PodFied
}

type dssResp struct {
	Total int                `json:"total"`
	Items []appsv1.DaemonSet `json:"items"`
}

func (d dsCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d dsCell) GetName() string {
	return d.Name
}

func (d *ds) toCells(dss []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(dss))
	for i := range dss {
		cells[i] = dsCell(dss[i])
	}
	return cells
}

func (d *ds) fromCells(cells []DataCell) []appsv1.DaemonSet {
	ds := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		ds[i] = appsv1.DaemonSet(cells[i].(dsCell))
	}
	return ds
}

func (d *ds) GetDsNum() (t []Dstotal, err error) {
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
		ds_list, err := K8s.Clientset.AppsV1().DaemonSets(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取damonset列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Dstotal{
			DsNum:     len(ds_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (d *ds) GetDs(filterName, namespace string, limit, page int) (*dssResp, error) {

	dsList, err := K8s.Clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取damonset列表失败" + err.Error())
		return nil, errors.New("获取damonset列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: d.toCells(dsList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	ds := d.fromCells(data.GenericDataList)

	return &dssResp{
		Total: total,
		Items: ds,
	}, nil
}

func (d *ds) GetDsDetails(name, namespace string) (ds *appsv1.DaemonSet, err error) {
	ds, err = K8s.Clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取damonset详情失败" + err.Error())
		return nil, errors.New("获取damonset详情失败" + err.Error())
	}
	return ds, nil
}

func (d *ds) DeleteDs(name, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除damonset失败" + err.Error())
		return errors.New("删除damonset失败" + err.Error())
	}
	return nil
}

func (d *ds) UpdateDs(name, namespace, content string) (err error) {
	var ds = &appsv1.DaemonSet{}
	err = json.Unmarshal([]byte(content), ds)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().DaemonSets(namespace).Update(context.TODO(), ds, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("damonset更新失败" + err.Error())
		return errors.New("damonset更新失败" + err.Error())
	}
	return nil
}

func (d *ds) CreateDs(data *DsFied) (err error) {

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

	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
	_, err = K8s.Clientset.AppsV1().DaemonSets(data.Namespace).Create(context.TODO(), ds, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建dsment失败" + err.Error())
		return errors.New("创建dsment失败" + err.Error())
	}
	return nil
}

func (d *ds) RestartDs(name, namespace string) (err error) {
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
	_, err = K8s.Clientset.AppsV1().DaemonSets(namespace).Patch(context.TODO(), name,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logrus.Error("重启Daemonset失败" + err.Error())
		return errors.New("重启Daemonset失败" + err.Error())
	}
	return nil
}
