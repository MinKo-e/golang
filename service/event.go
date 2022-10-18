package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var Event event

type event struct {
}

type EventsResp struct {
	Total int            `json:"total"`
	Items []corev1.Event `json:"items"`
}

type eventCell corev1.Event

func (e eventCell) GetCreation() time.Time {
	return e.CreationTimestamp.Time
}

func (e eventCell) GetName() string {
	return e.Name
}

func (e *event) toCells(events []corev1.Event) []DataCell {
	cells := make([]DataCell, len(events))
	for i := range events {
		cells[i] = eventCell(events[i])
	}
	return cells
}

func (e *event) fromCells(cells []DataCell) []corev1.Event {
	events := make([]corev1.Event, len(cells))
	for i := range cells {
		events[i] = corev1.Event(cells[i].(eventCell))
	}
	return events
}

func (e *event) GetEvents(filterName, namespace string, limit, page int) (*EventsResp, error) {

	eventList, err := K8s.Clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Error("获取event列表失败" + err.Error())
		return nil, errors.New("获取event列表失败" + err.Error())
	}
	//实例化结构体并填充字段
	selectoerQuery := dataSelector{GenericDataList: e.toCells(eventList.Items), dataSelectorQuery: &DataSelectorQuery{
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
	events := e.fromCells(data.GenericDataList)

	return &EventsResp{
		Total: total,
		Items: events,
	}, nil
}
