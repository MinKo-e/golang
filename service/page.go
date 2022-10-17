package service

// 分页处理函数，绑定到父结构体
func (d *dataSelector) Paging() *dataSelector {
	limit := d.dataSelectorQuery.PaginationQuery.Limit
	page := d.dataSelectorQuery.PaginationQuery.Page
	if limit <= 0 || page <= 0 {
		return d
	}
	//limit 10 page 1 end=10 start = 0
	endIndex := limit * page
	startIndex := endIndex - limit

	//判断endindex超出索引范围，如果超出，那么就默认数组长度为末尾索引
	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}
