package mq

type Producer interface {
	Send(key string, message []byte) error
}
