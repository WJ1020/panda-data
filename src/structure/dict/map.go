package dict

import "sync"

type Shard struct {
	m map[string] interface{}
	mutex sync.Mutex
}
type ConcurrentMap struct {
	table []*Shard
	count int32
}
