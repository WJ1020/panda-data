package reply

type SyntaxErrReply struct {
}

func (s *SyntaxErrReply) ToBytes() []byte {

	return []byte("-Err syntax error\r\n")
}

var wrongTypeErrBytes = []byte("-WRONG TYPE Operation against a key holding the wrong kind of value\r\n")

type WrongTypeErrReply struct{}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}
