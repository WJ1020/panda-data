package redis

type Client interface {
	Write([]byte) error
}
