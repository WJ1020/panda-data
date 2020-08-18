package db

var m = map[string]CmdFunc{
	"set": Set,
	"get": Get,
}

func resolverCmd(cmd string) (v CmdFunc, ok bool) {
	v, ok = m[cmd]
	return
}
