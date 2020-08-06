package db

import (
	"godis/src/interface/redis"
	"godis/src/redis/reply"
)

func Set(db *DB, args [][]byte) redis.Reply {

	return &reply.OkReply{}

}
