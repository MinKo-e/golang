package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/json"
	"k8sManager/config"
	"time"
)

var Pod pod

type pod struct{}

type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

// 定义pod结构体参数，用于接收前端数据
type PodFied struct {
	PName              string            `form:"pname" binding:"required" json:"pname"`
	Namespace          string            `form:"namespace" binding:"required" json:"namespace"`
	Labels             map[string]string `form:"labels"`
	CImage             string            `form:"cimage" binding:"required" json:"cimage"`
	CName              string            `form:"cname" binding:"required" json:"cname"`
	NodeSelector       map[string]string `form:"node_selector"`
	IImage             string            `form:"iimage"  json:"iimage"`
	IName              string            `form:"iname"  json:"iname"`
	ServiceAccountName string            `form:"service_account_name"`
	NodeName           string            `form:"node_name"`
	Cpu                string            `json:"cpu" form:"cpu"`
	Memory             string            `json:"memory" form:"memory"`
	ContainerPost      int32             `json:"container_post" form:"container_port"`
	HealthCheck        bool              `json:"health_check" form:"health_check"`
	HealthPath         string            `json:"health_path" form:"health_path"`
	HealthType         string            `json:"health_type" form:"health_type"`
	HealthHost         string            `json:"health_host" form:"health_host"`
	HealthExec         []string          `json:"health_exec" form:"health_exec"`
}

type Podtotal struct {
	PodNum    int
	Namespace string
}

// 定义corev1.pod数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

func (p *pod) GetPodNum() (t []Podtotal, err error) {
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
		pod_list, err := K8s.Clientset.CoreV1().Pods(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取Pod列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Podtotal{
			PodNum:    len(pod_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (p *pod) GetPods(filterName, namespace string, limit, page int) (*PodsResp, error) {

	podList, err := K8s.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pod列表失败" + err.Error())
		return nil, errors.New("获取Pod列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: p.toCells(podList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil
}

func (p *pod) GetPodDetails(name, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Pod详情失败" + err.Error())
		return nil, errors.New("获取Pod详情失败" + err.Error())
	}
	return pod, nil
}

func (p *pod) DeletePod(name, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Pod失败" + err.Error())
		return errors.New("删除Pod失败" + err.Error())
	}
	return nil
}

func (p *pod) CreatePod(podstruct *PodFied) (err error) {
	container := corev1.PodSpec{
		Containers: []corev1.Container{
			{Name: podstruct.CName, Image: podstruct.CImage, Ports: []corev1.ContainerPort{
				{Name: "http", Protocol: corev1.ProtocolTCP, ContainerPort: podstruct.ContainerPost},
			}},
		},
		NodeSelector:       podstruct.NodeSelector,
		ServiceAccountName: podstruct.ServiceAccountName,
		NodeName:           podstruct.NodeName,
	}
	if podstruct.IName != "" && podstruct.IImage != "" {
		container.InitContainers = []corev1.Container{
			{Name: podstruct.IName, Image: podstruct.IImage},
		}
	}
	if podstruct.HealthCheck {
		container.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler:        corev1.ProbeHandler{},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		if podstruct.HealthType == "tcp" {
			container.Containers[0].ReadinessProbe.ProbeHandler = corev1.ProbeHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: podstruct.ContainerPost,
					},
					Host: podstruct.HealthHost,
				},
			}
		} else if podstruct.HealthType == "http" {
			container.Containers[0].ReadinessProbe.ProbeHandler = corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: podstruct.ContainerPost,
					},
					Path: podstruct.HealthPath,
				},
			}
		} else if podstruct.HealthType == "exec" {
			container.Containers[0].ReadinessProbe.ProbeHandler = corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: podstruct.HealthExec,
				},
			}
		}
	}
	if podstruct.Cpu != "" && podstruct.Memory != "" {
		container.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(podstruct.Cpu),
			corev1.ResourceMemory: resource.MustParse(podstruct.Memory),
		}
		container.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(podstruct.Cpu),
			corev1.ResourceMemory: resource.MustParse(podstruct.Memory),
		}
	}
	option := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              podstruct.PName,
			Namespace:         podstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            podstruct.Labels,
		},
		Spec: container,
	}
	_, err = K8s.Clientset.CoreV1().Pods(podstruct.Namespace).Create(context.TODO(), &option, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建Pod失败" + err.Error())
		return errors.New("创建Pod失败" + err.Error())
	}
	return nil
}

func (p *pod) GetPodLog(name, namespace, containerName string) (log string, err error) {
	lineLimit := int64(config.PodLogLine)
	options := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}
	req := K8s.Clientset.CoreV1().Pods(namespace).GetLogs(name, options)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logrus.Error("获取Pod日志失败" + err.Error())
		return "", errors.New("获取Pod日志失败" + err.Error())
	}
	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {
			logrus.Error("获取Pod日志失败" + err.Error())
		}
	}(podLogs)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logrus.Error("复制Pod日志失败" + err.Error())
		return "", errors.New("复制Pod日志失败" + err.Error())
	}
	return buf.String(), nil
}

func (p *pod) UpdatePod(namespace, content string) (err error) {
	var pod = &corev1.Pod{}
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Pod更新失败" + err.Error())
		return errors.New("Pod更新失败" + err.Error())
	}
	return nil
}

func (p *pod) GetContainerName(name, namespace string) ([]string, error) {
	var containerList []string
	data, err := p.GetPodDetails(name, namespace)
	if err != nil {
		logrus.Error("获取Pod详情失败" + err.Error())
		return nil, errors.New("获取Pod详情失败" + err.Error())
	}
	for _, v := range data.Spec.Containers {
		containerList = append(containerList, v.Name)
	}
	return containerList, nil
}
