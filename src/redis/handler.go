package handler

import (
	"bufio"
	"context"
	DBImpl "godis/src/db"
	"godis/src/interface/db"
	"godis/src/lib/logger"
	"godis/src/lib/sync/atomic"
	"godis/src/redis/reply"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

var UnKnowErrReplyBytes = []byte("-ERR un know\r\n")

type Handler struct {
	activeConn sync.Map
	db         db.DB
	closing    atomic.BoolAtomic
}

func MakeHandler() *Handler {
	return &Handler{
		db: DBImpl.MakeDB(),
	}
}
func (h *Handler) closeClient(client *Client) {
	_ = client.Close()
	//TODO 通知数据库
	h.activeConn.Delete(client)

}

/**
用来解析redis协议
*/
func (h *Handler) Handler(ctx context.Context, conn net.Conn) {
	//数据库正在关闭中 拒绝所有新的客户端连接
	if h.closing.Get() {
		_ = conn.Close()
		return
	}
	client := MakeClient(conn)
	h.activeConn.Store(client, 1)

	reader := bufio.NewReader(conn)
	//下一个完整包的长度
	var fixLen int64 = 0
	var err error
	//本次读取的数据
	var msg []byte
	for {
		if fixLen == 0 {
			msg, err = reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					logger.Info("connection close ")
				} else {
					logger.Error(err)
				}

				h.closeClient(client)
				return
			}
			if len(msg) == 0 || msg[len(msg)-2] != '\r' {
				errReply := &reply.ProtocolErrReply{Msg: "invalid multi bulk length"}
				_, _ = client.conn.Write(errReply.ToBytes())
			}
		} else {
			//需要读取下一个完整的包+CRLF的长度
			msg := make([]byte, fixLen+2)
			_, err = io.ReadFull(reader, msg)
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					logger.Info("connection close")
				} else {
					logger.Error(err)
				}
				h.closeClient(client)
				return
			}
			if len(msg) == 0 || msg[len(msg)-2] != '\r' {
				errReply := &reply.ProtocolErrReply{Msg: "invalid multi bulk length"}
				_, _ = client.conn.Write(errReply.ToBytes())
			}
			fixLen = 0
		}

		if !client.uploading.Get() {
			//数组类型 * 开头
			if msg[0] == '*' {
				//获取该数组预期包含了多少个元素 放到 expectedLine 中
				expectedLine, err := strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
				if err != nil {
					_, _ = client.conn.Write(UnKnowErrReplyBytes)
					continue
				}
				client.waitingReplay.Add(1)
				client.uploading.Set(true)
				//有多少个元素
				client.expectedArgsCount = uint32(expectedLine)
				//已经收到了多少个
				client.receivedCount = 0
				//存放指令的数组
				client.args = make([][]byte, expectedLine)
			} else {
				//不是数组的类型 那么就是一行
				str := strings.TrimSuffix(string(msg), "\n")
				str = strings.TrimSuffix(str, "\r")
				//构建命令行为字符串数组
				strs := strings.Split(str, " ")
				//创建一个二进制的数组 将字符串数组转换为二进制数组
				args := make([][]byte, len(strs))
				for index, s := range strs {
					args[index] = []byte(s)
				}
				//TODO 发送给数据库执行 并且返回结果

				h.db.Exec(client, args)

				//返回给客户端结果
			}
		} else {
			//此时已经确定是数组的类型
			line := msg[0 : len(msg)-2]
			//看一下该数组的元素是不是大字符串的类型
			if line[0] == '$' {
				//大写字符串类型 查看该大字符串类型的长度为多少 下次应该读取该长度的数据
				fixLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
				if err != nil {
					errReply := &reply.ProtocolErrReply{Msg: err.Error()}
					_, _ = client.conn.Write(errReply.ToBytes())
				}
				if fixLen <= 0 {
					errReply := &reply.ProtocolErrReply{Msg: "无效的包长度"}
					_, _ = client.conn.Write(errReply.ToBytes())
				}

			} else {
				//不是大字符串就是整型或者简单字符串等类型
				client.args[client.receivedCount] = line
				client.receivedCount++
			}

			//查看是否已经接收完成
			if client.receivedCount == client.expectedArgsCount {
				client.uploading.Set(false)
				//TODO 发送给数据库执行并返回结果

				//本次需要接受的包长清零
				client.expectedArgsCount = 0
				//计数器清零
				client.receivedCount = 0
				client.args = nil
				client.waitingReplay.Done()
			}

		}

	}
}
func (h *Handler) Close() error {

	return nil
}
