package service

import "strings"

// 过滤函数，绑定到父结构体
func (d *dataSelector) Filter() *dataSelector {
	name := d.dataSelectorQuery.FilterQuery.Name
	if name == "" {
		return d
	}
	filterList := []DataCell{}
	for _, values := range d.GenericDataList {
		objName := values.GetName()
		if !strings.Contains(objName, name) {
			continue
		}
		filterList = append(filterList, values)
	}
	d.GenericDataList = filterList
	return d
}
