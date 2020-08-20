package reply

var (
	OK   = []byte("+OK\r\n")
	PONG = []byte("+PONG\r\n")
)

type OkReply struct {
}

type NullBulkReply struct {
}
type Pong struct {
}

func (p *Pong) ToBytes() []byte {
	return PONG
}

func (r *OkReply) ToBytes() []byte {
	return OK
}
func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkReplyBytes
}
