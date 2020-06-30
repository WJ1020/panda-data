package handler

import (
	"context"
	"godis/src/db"
	"godis/src/lib/sync/atomic"
	"net"
	"sync"
)

type Handler struct {
	activeConn sync.Map
	db db.DB
	closing atomic.BoolAtomic

}

func MakeHandler() *Handler{
	return &Handler{
		//TODO 初始化一个数据库
	}
}
func(h *Handler) closeClient(client *Client){
	_=client.Close()
	//TODO 通知数据库
	h.activeConn.Delete(client)

}
/**
	用来解析redis协议
 */
func(h *Handler) Handler(ctx context.Context,conn net.Conn){
	//数据库正在关闭中 拒绝所有新的客户端连接
	if h.closing.Get() {
		_=conn.Close()
		return
	}
	client:=MakeClient(conn)
	h.activeConn.Store(client,1)



}


