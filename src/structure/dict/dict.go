package dict

type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Put(key string, val interface{}) (result int)
	Remove(key string) (result int)
}
