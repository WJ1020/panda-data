package handler

import (
	"godis/src/lib/sync/atomic"
	"godis/src/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn              net.Conn
	waitingReplay     wait.Wait
	uploading         atomic.BoolAtomic
	expectedArgsCount uint32
	receivedCount     uint32
	args              [][]byte
	mu                sync.Mutex
	subs              map[string]bool
}

func (c *Client) Close() error {
	c.waitingReplay.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}
func MakeClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}
func (c *Client) Write(b []byte) error {
	if b == nil || len(b) == 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.conn.Write(b)
	return err
}
func (c *Client) SubsChannel(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.subs == nil {
		c.subs = make(map[string]bool)
	}
	c.subs[channel] = true
}
func (c *Client) SubsCount() int {
	if c.subs == nil {
		return 0
	}
	return len(c.subs)
}
func (c *Client) GetChannels() []string {
	if c.subs == nil {
		return make([]string, 0)
	}
	channels := make([]string, len(c.subs))
	i := 0
	for channel := range c.subs {
		channels[i] = channel
		i++
	}
	return channels
}
