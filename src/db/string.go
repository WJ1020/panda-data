package db

import (
	"panda-data/src/interface/redis"
	"panda-data/src/redis/reply"
	"strconv"
	"strings"
)

const (
	DEF = iota
	NX
	XX
)

func Set(db *DB, args [][]byte) redis.Reply {

	if len(args) < 2 {
		return reply.MakeErrRelay("Error comment set")
	}

	key := string(args[0])

	value := string(args[1])

	policy := DEF

	var ttl int64 = 0

	l := len(args)

	if l > 2 {
		//此时携带其他参数
		for i := 0; i < l; i++ {
			arg := strings.ToUpper(string(args[i]))
			if arg == "NX" {
				if policy == XX {
					return reply.MakeErrRelay("set commit error ")
				}
				policy = NX
			} else if arg == "XX" {
				if policy == NX {
					return reply.MakeErrRelay("set commit error")
				}
				policy = XX
			} else if arg == "EX" || arg == "PX" {
				if i+1 >= l {
					return &reply.StandardErrReply{}
				}
				s, err := strconv.ParseInt(string(args[i+1]), 10, 64)
				if err != nil {
					return &reply.SyntaxErrReply{}
				}
				if s < 0 {
					return reply.MakeErrRelay("Err expire time")
				}
				if arg == "EX" {
					ttl = s * 1000
				}

			}

		}

	}

	entity := &DataEntity{
		Data: value,
	}
	var result int
	switch policy {
	case DEF:
		result = db.Put(key, entity)
		//TODO 实现带过期实际的map
	case NX:
	case XX:

	}

	if ttl > 0 {
		//TODO
	}

	if result > 0 || policy == DEF {
		return &reply.OkReply{}
	}

	return &reply.OkReply{}

}
