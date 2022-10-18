package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type pvCell corev1.PersistentVolume

var Pv pv

type pv struct{}

type PvResp struct {
	Total int                       `json:"total"`
	Items []corev1.PersistentVolume `json:"items"`
}

type PvFied struct {
	Name                          string            `json:"name" binding:"required"`
	Label                         map[string]string `json:"label" binding:"required"`
	Storage                       string            `json:"storage"  binding:"required"`
	AccessModes                   string            `json:"access_modes" binding:"required"`
	StorageClass                  string            `json:"storage_class"`
	HostPath                      string            `json:"host_path"`
	NFS                           bool              `json:"nfs"`
	NFSServer                     string            `json:"nfs_server"`
	NFSPath                       string            `json:"nfs_path"`
	PersistentVolumeReclaimPolicy string            `json:"persistent_volume_reclaim_policy"`
	VolumeMode                    string            `json:"volume_mode"`
}

func (p pvCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p pvCell) GetName() string {
	return p.Name
}

func (p *pv) toCells(Pv []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(Pv))
	for i := range Pv {
		cells[i] = pvCell(Pv[i])
	}
	return cells
}

func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvm := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		pvm[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return pvm
}

func (p *pv) GetPvDetails(name string) (pvm *corev1.PersistentVolume, err error) {
	pvm, err = K8s.Clientset.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取pv详情失败" + err.Error())
		return nil, errors.New("获取pv详情失败" + err.Error())
	}
	return pvm, nil
}

func (p *pv) UpdatePv(content string) (err error) {
	var pvm = &corev1.PersistentVolume{}
	err = json.Unmarshal([]byte(content), pvm)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().PersistentVolumes().Update(context.TODO(), pvm, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Pv更新失败" + err.Error())
		return errors.New("Pv更新失败" + err.Error())
	}
	return nil
}

func (p *pv) CreatePv(data *PvFied) (err error) {
	var FileSystem = corev1.PersistentVolumeFilesystem
	var Block = corev1.PersistentVolumeBlock

	options := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Label,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: resource.MustParse(data.Storage),
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: data.HostPath,
					Type: nil,
				},
			},
			AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			StorageClassName:              "manual",
			VolumeMode:                    &FileSystem,
		},
	}
	if data.VolumeMode != "file_system" {
		options.Spec.VolumeMode = &Block
	}
	if data.AccessModes != "" {
		if data.AccessModes == "RWX" {
			options.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany}
		} else if data.AccessModes == "ROX" {
			options.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadOnlyMany}
		}

	}
	if data.StorageClass != "" {
		options.Spec.StorageClassName = data.StorageClass
	}
	if data.HostPath != "" {
		var Hostpath = &corev1.HostPathVolumeSource{
			Path: data.HostPath,
			Type: nil,
		}
		options.Spec.HostPath = Hostpath
	}
	if data.NFS && data.HostPath == "" {
		var nfs = &corev1.NFSVolumeSource{Server: data.NFSServer, Path: data.NFSPath, ReadOnly: false}
		options.Spec.PersistentVolumeSource.NFS = nfs
		options.Spec.PersistentVolumeSource.HostPath = nil
	}

	_, err = K8s.Clientset.CoreV1().PersistentVolumes().Create(context.TODO(), options, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建Pv失败" + err.Error())
		return errors.New("创建Pv失败" + err.Error())
	}
	return nil
}

func (p *pv) DeletePv(name string) (err error) {

	err = K8s.Clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除pv失败" + err.Error())
		return errors.New("删除pv失败" + err.Error())
	}
	return nil
}

func (p *pv) GetPv(filterName string, limit, page int) (*PvResp, error) {

	pvList, err := K8s.Clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pv列表失败" + err.Error())
		return nil, errors.New("获取Pv列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: p.toCells(pvList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	pv := p.fromCells(data.GenericDataList)

	return &PvResp{
		Total: total,
		Items: pv,
	}, nil
}

func (p *pv) GetPvNum() (t int, err error) {

	pv_list, err := K8s.Clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Pod列表失败" + err.Error())
		return 0, err

	}
	t = len(pv_list.Items)

	return t, nil
}

func (p *pv) GetPvBind(name string) (n string, err error) {

	pv, err := K8s.Clientset.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Pod列表失败" + err.Error())
		return "", err

	}
	n = ""
	if n = pv.Spec.ClaimRef.Name; n == "" {
		return "", nil
	}

	return n, nil
}
