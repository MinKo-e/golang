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

var Secret secret

type secret struct{}

type SecretsResp struct {
	Total int             `json:"total"`
	Items []corev1.Secret `json:"items"`
}

// 定义secret结构体参数，用于接收前端数据
type SecretFied struct {
	Name      string            `form:"name" binding:"required" json:"name"`
	Namespace string            `form:"namespace" binding:"required" json:"namespace"`
	Label     map[string]string `form:"label" json:"labels"`
	Data      map[string]string `form:"data" json:"data"`
	Type      string            `json:"type" form:"type"`
}

type Secrettotal struct {
	SecretNum int
	Namespace string
}

var Secretfied SecretFied

var DaTa map[string][]byte

func (s *SecretFied) toBytes() {
	for k, v := range s.Data {
		DaTa[k] = []byte(v)
	}
}

// 定义corev1.secret数据类型，实现Datacell接口，也就是实现了datacell数据类型，实现了dataselector结构体GenericDataList字段的数据属性
type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}

func (s *secret) toCells(secrets []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(secrets))
	for i := range secrets {
		cells[i] = secretCell(secrets[i])
	}
	return cells
}

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

func (s *secret) GetSecretNum() (t []Secrettotal, err error) {
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
		secret_list, err := K8s.Clientset.CoreV1().Secrets(v).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Error("获取Secret列表失败" + err.Error())
			return nil, err

		}
		t = append(t, Secrettotal{
			SecretNum: len(secret_list.Items),
			Namespace: v,
		})
	}

	return t, nil
}

func (s *secret) GetSecrets(filterName, namespace string, limit, page int) (*SecretsResp, error) {

	secretList, err := K8s.Clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取Secret列表失败" + err.Error())
		return nil, errors.New("获取Secret列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: s.toCells(secretList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	secrets := s.fromCells(data.GenericDataList)

	return &SecretsResp{
		Total: total,
		Items: secrets,
	}, nil
}

func (s *secret) GetSecretDetails(name, namespace string) (secret *corev1.Secret, err error) {
	secret, err = K8s.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logrus.Error("获取Secret详情失败" + err.Error())
		return nil, errors.New("获取Secret详情失败" + err.Error())
	}
	return secret, nil
}

func (s *secret) DeleteSecret(name, namespace string) (err error) {
	err = K8s.Clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error("删除Secret失败" + err.Error())
		return errors.New("删除Secret失败" + err.Error())
	}
	return nil
}

func (s *secret) CreateSecret(secretstruct *SecretFied) (err error) {
	option := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:              secretstruct.Name,
			Namespace:         secretstruct.Namespace,
			CreationTimestamp: metav1.Time{},
			Labels:            secretstruct.Label,
		},
		Data: DaTa,
		Type: corev1.SecretTypeServiceAccountToken,
	}
	if secretstruct.Type != "" {
		if secretstruct.Type == "registry" {
			option.Type = corev1.SecretTypeDockercfg
		} else if secretstruct.Type == "tls" {
			option.Type = corev1.SecretTypeTLS
		} else if secretstruct.Type == "basic" {
			option.Type = corev1.SecretTypeBasicAuth
		} else if secretstruct.Type == "opaque" {
			option.Type = corev1.SecretTypeOpaque
		}
	}
	_, err = K8s.Clientset.CoreV1().Secrets(secretstruct.Namespace).Create(context.TODO(), option, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("创建Secret失败" + err.Error())
		return errors.New("创建Secret失败" + err.Error())
	}
	return nil
}

func (s *secret) UpdateSecret(namespace, content string) (err error) {
	var secret = &corev1.Secret{}
	err = json.Unmarshal([]byte(content), secret)
	if err != nil {
		logrus.Error("Json反序列化失败" + err.Error())
		return errors.New("Json反序列化失败" + err.Error())
	}
	_, err = K8s.Clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("Secret更新失败" + err.Error())
		return errors.New("Secret更新失败" + err.Error())
	}
	return nil
}

func (s *secret) GetSecretData(name, namespace string) (map[string]string, error) {
	var rdata map[string]string
	data, err := s.GetSecretDetails(name, namespace)
	if err != nil {
		logrus.Error("获取Secret详情失败" + err.Error())
		return nil, errors.New("获取Secret详情失败" + err.Error())
	}
	for k, v := range data.Data {
		rdata[k] = string(v)
	}
	return rdata, nil
}
