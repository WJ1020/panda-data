package db

import (
	"panda-data/src/interface/redis"
	"panda-data/src/redis/reply"
)

func Set(db *DB, args [][]byte) redis.Reply {

	return &reply.OkReply{}

}
