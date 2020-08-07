package db

import (
	"panda-data/src/interface/redis"
)

type DB interface {
	Exec(client redis.Client, args [][]byte) redis.Reply
}
