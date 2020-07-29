package db

import (
	"godis/src/interface/client"
	"godis/src/interface/redis"
)

type DB interface {
	Exec(client client.Client, args [][]byte) redis.Reply
}
