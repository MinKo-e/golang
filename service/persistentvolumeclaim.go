package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"time"
)

var Pvc pvc

type pvc struct{}

type PvcsResp struct {
	Total int                            `json:"total"`
	Items []corev1.PersistentVolumeClaim `json:"items"`
}

// 定义pvc结构体参数，用于接收前端数据
type PvcFied struct {
	Name                          string            `json:"name" binding:"required"`
	Namespace                     string            `json:"namespace" binding:"required"`
	Label                         map[string]string `json:"label" binding:"required"`
	Storage                       string            `json:"storage"  binding:"required"`
	AccessModes                   string            `json:"access_modes" binding:"required" `
	StorageClass                  string            `json:"storage_class"`
	PersistentVolumeReclaimPolicy string            `json:"persistent_volume_reclaim_policy"`
	VolumeMode                    string            `json:"volume_mode"`
	DataSourceName                string            `json:"data_source_name"`
}

type Pvctotal struct {
	PvcNum    int
	Namespace string
}

// 定义corev1.pvc数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type pvcCell corev1.PersistentVolumeClaim

func (p pvcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvcCell) GetName() string {
	return p.Name
}

func (p *pvc) toCells(pvcs []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(pvcs))
	for i := range pvcs {
		cells[i] = pvcCell(pvcs[i])
	}
	return cells
}

func (p *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	pvcs := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		pvcs[i] = corev1.PersistentVolumeClaim(cells[i].(pvcCell))
	}
	return pvcs
}

func (p *pvc) GetPvcNum() (t []Pvctotal, err error) {
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
		pvc_list, err := K8s.Clientset.CoreV1().PersistentVolumeClaims(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取Pvc列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Pvctotal{
			PvcNum:    len(pvc_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (p *pvc) GetPvcs(filterName, namespace string, limit, page int) (*PvcsResp, error) {

	pvcList, err := K8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pvc列表失败" + err.Error())
		return nil, errors.New("获取Pvc列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: p.toCells(pvcList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	pvcs := p.fromCells(data.GenericDataList)

	return &PvcsResp{
		Total: total,
		Items: pvcs,
	}, nil
}

func (p *pvc) GetPvcDetails(name, namespace string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = K8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Pvc详情失败" + err.Error())
		return nil, errors.New("获取Pvc详情失败" + err.Error())
	}
	return pvc, nil
}

func (p *pvc) DeletePvc(name, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Pvc失败" + err.Error())
		return errors.New("删除Pvc失败" + err.Error())
	}
	return nil
}

func (p *pvc) CreatePvc(pvcstruct *PvcFied) (err error) {
	var FileSystem = corev1.PersistentVolumeFilesystem
	var Block = corev1.PersistentVolumeBlock
	var SC = "manual"
	DataSource := &corev1.TypedLocalObjectReference{
		APIGroup: nil,
		Kind:     "PersistentVolumeClaim",
		Name:     pvcstruct.DataSourceName,
	}
	option := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:              pvcstruct.Name,
			Namespace:         pvcstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            pvcstruct.Label,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(pvcstruct.Storage),
				},
			},

			StorageClassName: &SC,
			VolumeMode:       &FileSystem,
		},
	}
	if pvcstruct.VolumeMode != "" {
		if pvcstruct.VolumeMode != "file_system" {
			option.Spec.VolumeMode = &Block
		}
	}
	if pvcstruct.AccessModes != "" {
		if pvcstruct.AccessModes == "RWX" {
			option.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany}
		} else if pvcstruct.AccessModes == "ROX" {
			option.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadOnlyMany}
		}

	}
	if pvcstruct.StorageClass != "" {
		option.Spec.StorageClassName = &pvcstruct.StorageClass
	}
	if pvcstruct.DataSourceName != "" {
		option.Spec.DataSource = DataSource
	}
	_, err = K8s.Clientset.CoreV1().PersistentVolumeClaims(pvcstruct.Namespace).Create(context.TODO(), option, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建Pvc失败" + err.Error())
		return errors.New("创建Pvc失败" + err.Error())
	}
	return nil
}

func (p *pvc) UpdatePvc(namespace, content string) (err error) {
	var pvc = &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal([]byte(content), pvc)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Pvc更新失败" + err.Error())
		return errors.New("Pvc更新失败" + err.Error())
	}
	return nil
}
