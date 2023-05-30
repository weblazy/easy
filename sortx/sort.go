package sortx

// 按照接口中的某个字段排序
type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

type Sort struct {
	Obj  interface{} `json:"obj"`
	Sort int64       `json:"sort"`
}
type SortList struct {
	List  []Sort
	Order Order
}

func NewSortList(order Order) *SortList {
	return &SortList{
		List:  make([]Sort, 0),
		Order: order,
	}
}

func (list *SortList) Len() int {
	return len(list.List)
}

func (list *SortList) Less(i, j int) bool {
	if list.Order == ASC {
		return list.List[i].Sort <= list.List[j].Sort
	}
	return list.List[i].Sort >= list.List[j].Sort
}

func (list *SortList) Swap(i, j int) {
	temp := list.List[i]
	list.List[i] = list.List[j]
	list.List[j] = temp
}
