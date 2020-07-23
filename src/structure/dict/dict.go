package dict

type Dict interface {
	Get(key string) (val interface{}, exists bool)
}
