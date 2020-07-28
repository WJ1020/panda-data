package dict


/**
	哈希表的节点
 */
type dictEntry struct {
	key string
	val interface{}
	next *dictEntry
}

/**
	哈希表
 */
type dictHt struct {

	/**
	哈希表
	 */
	table []dictEntry
	/**
	哈希表大小
	 */
	size uint64
	/**
	 size-1
	 */
	sizeMask uint64
	/**
	已经使用的节点数量
	 */
	used uint64


}

type dict struct {

	privateData interface{}

	ht[2] dictHt

	rehashIndex int

	iterators int

	
}

func (ht *dictHt) makeDictHt() *dictHt  {
	ht.table =nil
	ht.size=0
	ht.sizeMask=0
	ht.used=0
	return ht
}


func (d * dict) makeDict(privateData interface{}) *dict{
	d.privateData=privateData
	d.rehashIndex=-1
	d.iterators=0
	d.ht[0].makeDictHt()
	d.ht[1].makeDictHt()
	return d
}


