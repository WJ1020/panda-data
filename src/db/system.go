package db

import (
	"panda-data/src/interface/redis"
	"panda-data/src/redis/reply"
	"strconv"
)

func Ping(db *DB, args [][]byte) redis.Reply {
	return &reply.Pong{}
}

//TODO 待完善
func Select(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return reply.MakeErrRelay("Error select")
	}
	s := string(args[0])
	v, e := strconv.ParseInt(s, 10, 32)
	if e != nil {
		return reply.MakeErrRelay("Error select")
	}
	if v < 3 {
		return &reply.OkReply{}
	} else {
		return reply.MakeErrRelay("Index of ")
	}

}
