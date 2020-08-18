package reply

type OkReply struct {
}

type NullBulkReply struct {
}

func (r *OkReply) ToBytes() []byte {
	return []byte("+OK\r\n")
}
func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkReplyBytes
}
