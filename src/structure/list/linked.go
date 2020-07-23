package list

/**
链表实现
*/
type LinkList struct {
	first *node
	last  *node
	size  int
}

/**
数据节点
*/
type node struct {
	data interface{}
	prev *node
	next *node
}

func (list *LinkList) Add(data interface{}) {
	if list == nil {
		panic("list is nil")
	}
	//新建一个节点
	n := &node{
		data: data,
	}
	if list.last == nil {
		list.first = n
		list.last = n
	} else {
		n.prev = list.last
		list.last.next = n
		list.last = n
	}
	list.size++
}
