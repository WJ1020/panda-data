package db

import (
	"panda-data/src/interface/redis"
	"panda-data/src/redis/reply"
)

func Del(db *DB, args [][]byte) redis.Reply {
	if len(args) == 0 {
		return reply.MakeErrRelay("Error wrong del command")
	}
	//TODO 针对KEY的批量加锁和批量移除

	return nil

}
