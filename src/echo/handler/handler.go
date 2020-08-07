package handlerEcho

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"panda-data/src/lib/sync/atomic"
	"panda-data/src/lib/sync/wait"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.BoolAtomic
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type Client struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *Client) Close() error {
	c.Waiting.WaitWithTimeout(8 * time.Second)
	err := c.Conn.Close()
	log.Println(fmt.Sprintf("CLOSE CONNECTion %v", err))
	return nil
}
func (h *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		err := conn.Close()
		log.Println(fmt.Sprintf("CLOSE connection %v", err))
	}
	client := &Client{
		Conn: conn,
	}
	h.activeConn.Store(client, 1)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection close")
			} else {
				log.Fatal(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, er := conn.Write(b)
		if er != nil {
			log.Fatal(er)
		}
		client.Waiting.Done()
	}
}
func (h *EchoHandler) Close() error {
	log.Println("handler shut ing down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*Client)
		_ = client.Close()
		return true
	})
	return nil
}
