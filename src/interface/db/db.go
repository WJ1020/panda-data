package db

import (
	"godis/src/interface/redis"
)

type DB interface {
	Exec(client redis.Client, args [][]byte) redis.Reply
}
