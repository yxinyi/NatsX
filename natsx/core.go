package natsx

import (
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"log"
	"reflect"
)

var Obj = New()

func New() *Struct {
	tmpObj := &Struct{
		CallbackPoll: make(map[Enum]reflect.Value),
		RedisCli:     redis.NewClient(&redis.Options{Addr: ":6379", DB: 15}),
	}
	return tmpObj
}

var connTimes = 3

func (n *Struct) Connection(url string) {
	n.ConnUrl = url
	for idx := 0; idx < connTimes; idx++ {
		conn, err := nats.Connect(n.ConnUrl)
		if err != nil {
			log.Printf("[natsx:Connection][%v][%v]", err.Error(), n.ConnUrl)

			continue
		}
		n.conn = conn
		break
	}
}

type Struct struct {
	conn         *nats.Conn
	ConnUrl      string
	CallbackPoll map[Enum]reflect.Value
	RedisCli     *redis.Client
}
