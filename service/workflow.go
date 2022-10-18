package service

import (
	"fmt"
	"k8sManager/dao"
	"k8sManager/model"
)

var Workflow workflow

type workflow struct {
}

type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int                    `json:"replicas"`
	Labels        map[string]string      `json:"labels"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
	ServiceType   string                 `json:"service_type"`
	Cpu           string                 `json:"cpu" form:"cpu"`
	Memory        string                 `json:"memory" form:"memory"`
	ContainerPost int32                  `json:"container_post" form:"container_port"`
	HealthCheck   bool                   `json:"health_check" form:"health_check"`
	HealthPath    string                 `json:"health_path" form:"health_path"`
	HealthType    string                 `json:"health_type" form:"health_type"`
	HealthHost    string                 `json:"health_host" form:"health_host"`
	HealthExec    []string               `json:"health_exec" form:"health_exec"`
	Image         string                 `json:"image"`
	Annotations   map[string]string      `json:"annotations"`
	Type          string                 `json:"type"`
	NodePort      int                    `form:"node_port" json:"node_port"`
	Port          int                    `form:"port" json:"port"`
}

func (w *workflow) GetList(name string, limit, page int) (*dao.WorkResp, error) {
	data, err := dao.Workflow.GetList(name, limit, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) GetByID(id int) (*model.Workflow, error) {
	data, err := dao.Workflow.GetByID(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) DeleteWorkFlow(id int) error {
	data, err := dao.Workflow.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = w.Delete(data)
	if err != nil {
		return err
	}
	return nil
}

func (w *workflow) Delete(data *model.Workflow) error {
	err := Deploy.DeleteDeploy(data.Name, data.Namespace)
	if err != nil {
		return err
	}
	err = NetworkSvc.DeleteNetworkSvc(getServiceName(data.Name), data.Namespace)
	if err != nil {
		return err
	}
	if data.ServiceType == "Ingress" {
		err = NetworkIngress.DeleteNetworkIngress(getIngressName(data.Name), data.Namespace)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *workflow) CreateWorkFlow(data *WorkflowCreate) error {
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}
	workflow := model.Workflow{
		Name:        data.Name,
		Namespace:   data.Namespace,
		Replicas:    int32(data.Replicas),
		Deployment:  data.Name,
		Service:     getServiceName(data.Name),
		Ingress:     ingressName,
		ServiceType: data.ServiceType,
	}
	if err := dao.Workflow.Create(&workflow); err != nil {
		return err
	}
	return nil

}

func (w *workflow) Create(wc *WorkflowCreate) error {

	deploy := &DeployFied{
		Name:     wc.Name,
		Replicas: int32(wc.Replicas),
		Label:    wc.Labels,
		PodFied: PodFied{
			PName:         wc.Name,
			Namespace:     wc.Namespace,
			Labels:        wc.Labels,
			CImage:        wc.Image,
			CName:         wc.Name,
			NodeSelector:  nil,
			Cpu:           wc.Cpu,
			Memory:        wc.Memory,
			ContainerPost: wc.ContainerPost,
			HealthCheck:   wc.HealthCheck,
			HealthPath:    wc.HealthPath,
			HealthType:    wc.HealthType,
			HealthHost:    wc.HealthHost,
			HealthExec:    wc.HealthExec,
		},
	}
	err := Deploy.CreateDeploy(deploy)
	if err != nil {
		err = Deploy.DeleteDeploy(wc.Name, wc.Namespace)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err
	}
	service := &NetworkSvcFied{
		Name:       getServiceName(wc.Name),
		Namespace:  wc.Namespace,
		Labels:     wc.Labels,
		Type:       wc.Type,
		Selector:   wc.Labels,
		NodePort:   wc.NodePort,
		TargetPort: int(wc.ContainerPost),
		Port:       wc.Port,
	}
	err = NetworkSvc.CreateNetworkSvc(service)
	if err != nil {
		err = NetworkSvc.DeleteNetworkSvc(getServiceName(wc.Name), wc.Namespace)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err
	}
	if wc.ServiceType == "Ingress" {
		ingress := &NetworkIngressFied{
			Name:        getIngressName(wc.Name),
			Namespace:   wc.Namespace,
			Labels:      wc.Labels,
			Annotations: wc.Annotations,
			Hosts:       wc.Hosts,
		}
		err := NetworkIngress.CreateNetworkIngress(ingress)
		if err != nil {
			err = NetworkIngress.DeleteNetworkIngress(getIngressName(wc.Name), wc.Namespace)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return err
		}
		return err
	}
	err = w.CreateWorkFlow(wc)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func getIngressName(name string) string {
	return name + "-ing"
}
func getServiceName(name string) string {
	return name + "-svc"
}
