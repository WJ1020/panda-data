package reply

import "strconv"

var(
	nullBulkReplyBytes=[]byte("$-1")
	CRLF ="\r\n"
)

type BulkReply struct {
	Arg []byte
}

func MakeBulkReply(arg []byte) *BulkReply{
	return &BulkReply{
		Arg: arg,
	}
}
func (r *BulkReply) ToBytes() []byte{
		if len(r.Arg)==0{
			return nullBulkReplyBytes
		}
		s:="$"+strconv.Itoa(len(r.Arg))+CRLF+string(r.Arg)+CRLF
		return []byte(s)
}

type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply{
	return &MultiBulkReply{
		Args: args,
	}
}

func(r *MultiBulkReply) ToBytes() []byte {
	argLen := len(r.Args)
	res := "*" + strconv.Itoa(argLen) + CRLF
	for _, arg := range r.Args {
		if arg == nil {
			res += "$-1" + CRLF
		} else {
			res += "$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF
		}
	}
	return []byte(res)
}

type StatusReply struct {
	Status string
}
func MakeStatusReply(status string) *StatusReply{
	return &StatusReply{
		Status: status,
	}
}

type IntReply struct {
	Code int64
}
func MakeIntReply(code int64) *IntReply{
	return &IntReply{
		Code: code,
	}
}
func(r *IntReply) ToBytes() []byte{
	s:=":"+strconv.FormatInt(r.Code,10)+CRLF
	return []byte(s)
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandardErrReply struct {
	Status string
}
func MakeErrRelay(status string) *StandardErrReply{
	return &StandardErrReply{
		Status: status,
	}
}

func (r *StandardErrReply) ToBytes() []byte{
	return []byte("-"+r.Status+"\r\n")
}

func(r *StandardErrReply) Error() string{
	return r.Status
}

type ProtocolErrReply struct {
	Msg string
}
func (p *ProtocolErrReply) ToBytes() []byte{
	return []byte(p.Msg)
}