package handler

import (
	"bufio"
	"context"
	"godis/src/db"
	"godis/src/lib/logger"
	"godis/src/lib/sync/atomic"
	"godis/src/redis/reply"
	"io"
	"net"
	"strconv"
	"sync"
)

var UnKnowErrReplyBytes = []byte("-ERR un know\r\n")

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

	reader := bufio.NewReader(conn)

	var fixLen int64=0
	var err error
	var msg []byte
	for  {
		if fixLen ==0{
			msg,err=reader.ReadBytes('\n')
			if err!=nil {
				if err==io.EOF || err == io.ErrUnexpectedEOF {
					logger.Info("connection close ")
				}else {
					logger.Error(err)
				}

				h.closeClient(client)
				return
			}
			if len(msg) ==0 || msg[len(msg)-2]!='\r' {
				errReply:=&reply.ProtocolErrReply{Msg: "invalid multi bulk length"}
				_,_=client.conn.Write(errReply.ToBytes())
			}
		}else{
			msg:=make([]byte,fixLen+2)
			_,err=io.ReadFull(reader,msg)
			if err !=nil {
				if err == io.EOF||err == io.ErrUnexpectedEOF {
					logger.Info("connection close")
				}else {
					logger.Error(err)
				}
				h.closeClient(client)
				return
			}
			if len(msg) ==0 || msg[len(msg)-2]!='\r' {
				errReply:=&reply.ProtocolErrReply{Msg: "invalid multi bulk length"}
				_,_=client.conn.Write(errReply.ToBytes())
			}
			fixLen=0
		}

		if !client.uploading.Get(){
			if msg[0] == '*' {
				expectedLine,err := strconv.ParseUint(string(msg[1:len(msg)-2]),10,32)
				if err!=nil{
					_,_=client.conn.Write(UnKnowErrReplyBytes)
					continue
				}
				client.waitingReplay.Add(1)
				client.uploading.Set(true)
				client.expectedArgsCount=uint32(expectedLine)
				client.receivedCount=0
				client.args=make([][]byte,expectedLine)
			}else {

			}
		}
	}
}


