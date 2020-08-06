package reply

type OkReply struct {
}

func (r *OkReply) ToBytes() []byte {
	return []byte("+OK\r\n")
}
