package db

import (
	"godis/src/interface/client"
	"godis/src/interface/redis"
	"godis/src/lib/logger"
	"godis/src/structure/dict"
	"strings"
)

const (
	dataDictSize = 1 << 16
)

type DB struct {
	Data dict.Dict
}

func MakeDB() *DB {
	db := &DB{
		Data: dict.MakeConcurrent(dataDictSize),
	}
	return db
}

func (db *DB) Exec(c *client.Client, args [][]byte) (result redis.Reply) {
	cmd := strings.ToLower(string(args[0]))
	//执行的命令为
	logger.Info(cmd)
	return nil
}
