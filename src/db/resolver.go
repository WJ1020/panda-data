package db

var m = map[string]CmdFunc{}

func resolverCmd(cmd string) (v CmdFunc, ok bool) {
	v, ok = m[cmd]
	return
}
