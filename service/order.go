package service

import "sort"

// 实现排序方法
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return a.Before(b)
}

// 调用sort包实现对原始数据的排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}
