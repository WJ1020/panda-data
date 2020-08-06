package reply

type SyntaxErrReply struct {
}

func (s *SyntaxErrReply) ToBytes() []byte {

	return []byte("-Err syntax error\r\n")
}
