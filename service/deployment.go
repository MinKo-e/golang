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

type deployCell appsv1.Deployment

var Deploy deploy

type deploy struct{}

type Deploytotal struct {
	DeployNum int
	Namespace string
}

type DeployFied struct {
	Name     string            `json:"name"`
	Replicas int32             `json:"replicas"`
	Label    map[string]string `json:"label"`
	PodFied
}

type DeploysResp struct {
	Total int                 `json:"total"`
	Items []appsv1.Deployment `json:"items"`
}

func (d deployCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deployCell) GetName() string {
	return d.Name
}

func (d *deploy) toCells(Deploys []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(Deploys))
	for i := range Deploys {
		cells[i] = deployCell(Deploys[i])
	}
	return cells
}

func (d *deploy) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deployCell))
	}
	return deployments
}

func (d *deploy) GetDeployNum() (t []Deploytotal, err error) {
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
		deploy_list, err := K8s.Clientset.AppsV1().Deployments(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取Deploy列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Deploytotal{
			DeployNum: len(deploy_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (d *deploy) GetDeploys(filterName, namespace string, limit, page int) (*DeploysResp, error) {

	DeployList, err := K8s.Clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Deploy列表失败" + err.Error())
		return nil, errors.New("获取Deploy列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: d.toCells(DeployList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	deploys := d.fromCells(data.GenericDataList)

	return &DeploysResp{
		Total: total,
		Items: deploys,
	}, nil
}

func (d *deploy) GetDeployDetails(name, namespace string) (deploy *appsv1.Deployment, err error) {
	deploy, err = K8s.Clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Deploy详情失败" + err.Error())
		return nil, errors.New("获取Deploy详情失败" + err.Error())
	}
	return deploy, nil
}

func (d *deploy) DeleteDeploy(name, namespace string) (err error) {
	err = K8s.Clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Deploy失败" + err.Error())
		return errors.New("删除Deploy失败" + err.Error())
	}
	return nil
}

func (d *deploy) ScaleDeploys(name, namespace string, replicas int) (replica int32, err error) {

	scale, err := K8s.Clientset.AppsV1().Deployments(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取副本数失败" + err.Error())
		return 0, errors.New("获取副本数失败" + err.Error())
	}
	scale.Spec.Replicas = int32(replicas)
	newscale, err := K8s.Clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("更新副本数失败" + err.Error())
		return 0, errors.New("更新副本数失败" + err.Error())
	}
	return newscale.Spec.Replicas, nil
}

func (d *deploy) RestartDeploy(name, namespace string) (err error) {
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
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Patch(context.TODO(), name,
		"application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logrus.Error("重启Deployment失败" + err.Error())
		return errors.New("重启Deployment失败" + err.Error())
	}
	return nil
}

func (d *deploy) UpdateDeploy(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Deploy更新失败" + err.Error())
		return errors.New("Deploy更新失败" + err.Error())
	}
	return nil
}

func (d *deploy) CreateDeploy(data *DeployFied) (err error) {

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

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
	_, err = K8s.Clientset.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建deployment失败" + err.Error())
		return errors.New("创建deployment失败" + err.Error())
	}
	return nil
}
