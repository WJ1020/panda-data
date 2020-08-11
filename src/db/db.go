package db

import (
	"panda-data/src/interface/redis"
	"panda-data/src/redis/reply"
	"panda-data/src/structure/dict"
	"strings"
	"sync"
)

const (
	dataDictSize = 1 << 16
)

type DataEntity struct {
	Data interface{}
}

type DB struct {
	Data dict.Dict

	countDownLatch sync.WaitGroup
}

func MakeDB() *DB {
	db := &DB{
		Data: dict.MakeConcurrent(dataDictSize),
	}
	return db
}

type CmdFunc func(db *DB, args [][]byte) redis.Reply

func (db *DB) Exec(c redis.Client, args [][]byte) (result redis.Reply) {
	cmd := strings.ToLower(string(args[0]))
	cmdFunc, ok := resolverCmd(cmd)
	if !ok {
		//该命令暂时不支持
		return reply.MakeErrRelay("ERR unknown command: " + cmd)
	}
	if len(args) > 1 {
		result = cmdFunc(db, args[1:])
	} else {
		result = cmdFunc(db, [][]byte{})
	}
	return
}

func (db *DB) Put(key string, val *DataEntity) int {
	db.countDownLatch.Wait()
	return db.Data.Put(key, val)
}

func (db *DB) Get(key string) (*DataEntity, bool) {
	db.countDownLatch.Wait()
	raw, ok := db.Data.Get(key)
	if !ok {
		return nil, false
	}
	//TODO 判断key是否已经过期
	//interface可以通过这种方式强制转换
	res, _ := raw.(*DataEntity)
	return res, true
}
