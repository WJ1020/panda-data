package atomic

import "sync/atomic"

type BoolAtomic uint32
func(b *BoolAtomic) Get() bool{
	return atomic.LoadUint32((*uint32)(b))!=0
}
func(b *BoolAtomic) Set(v bool){
	if v{
		atomic.StoreUint32((*uint32)(b),1)
	}else {
		atomic.StoreUint32((*uint32)(b),0)
	}
}