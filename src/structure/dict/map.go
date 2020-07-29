package dict

import (
	"math"
	"sync"
	"sync/atomic"
)

const prime32 = uint32(16777619)

type Shard struct {
	m     map[string]interface{}
	mutex sync.Mutex
}
type ConcurrentMap struct {
	table []*Shard
	count int32
}

/**
找到2的整数次方
*/
func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	} else {
		return int(n + 1)
	}
}
func MakeConcurrent(shardCount int) *ConcurrentMap {
	shardCount = computeCapacity(shardCount)
	table := make([]*Shard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &Shard{
			m: make(map[string]interface{}),
		}
	}
	d := &ConcurrentMap{
		count: 0,
		table: table,
	}
	return d
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

//判断该哈希值应该在第几个槽位
func (dict *ConcurrentMap) spread(hashCode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}
	tableSize := uint32(len(dict.table))
	return (tableSize - 1) & hashCode
}

/**
根据槽位index获取槽位的指针
*/
func (dict *ConcurrentMap) getShard(index uint32) *Shard {
	if dict == nil {
		panic("dict is nil")
	}
	return dict.table[index]
}

func (dict *ConcurrentMap) Get(key string) (val interface{}, ext bool) {
	if dict == nil {
		panic("dict is nil")
	}
	//获取哈希值
	hash := fnv32(key)
	//获取槽位
	index := dict.spread(hash)
	//获取具体的分片
	shard := dict.table[index]
	//加锁
	shard.mutex.Lock()
	//函数结束时解锁
	defer shard.mutex.Unlock()
	//获取值
	val, ext = shard.m[key]
	//返回结果
	return val, ext
}

func (dict *ConcurrentMap) Put(key string, val interface{}) (res int) {
	if dict == nil {
		panic("dict is nil")
	}
	//获取哈希值
	hash := fnv32(key)
	//获取槽位
	index := dict.spread(hash)
	//获取具体的分片
	shard := dict.table[index]
	//加锁
	shard.mutex.Lock()
	//函数结束时解锁
	defer shard.mutex.Unlock()

	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return 0
	} else {
		shard.m[key] = val
		dict.addCount()
		return 1
	}
}

func (dict *ConcurrentMap) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}
