package sortx

//按照接口中的某个字段排序
type Sort struct {
	Obj  interface{} `json:"obj"`
	Sort int64       `json:"sort"`
}
type SortList []Sort

func (list SortList) Len() int {
	return len(list)
}

func (list SortList) Less(i, j int) bool {
	return list[i].Sort <= list[j].Sort
}

func (list SortList) Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
}
