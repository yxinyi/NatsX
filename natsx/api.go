package natsx

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"log"
	"reflect"
	"time"
)

var syncTimes = 3

func SyncMsg(enumType Enum, msg interface{}, retParam interface{}) interface{} {
	for idx := 0; idx < syncTimes; idx++ {
		coreMsgStream := ToCoreMsgStream(msg)
		backMsg, reqErr := Obj.conn.Request(enumType.ToString(), coreMsgStream, time.Second)
		if reqErr != nil {
			log.Printf("[SyncMsg][%v]", reqErr.Error())
			continue
		}
		if retParam == nil {
			return nil
		}
		if len(backMsg.Data) == 0 {
			return nil
		}

		retRefType := reflect.TypeOf(retParam)
		retRefVal := reflect.New(retRefType)
		jsoniter.Unmarshal(backMsg.Data, retRefVal.Interface())
		return retRefVal
	}
	return nil
}

func RegisterMsg(enumType Enum, cb interface{}) {
	cbVal := reflect.ValueOf(cb)
	_, exists := Obj.CallbackPoll[enumType]
	if exists {
		panic(fmt.Sprintf("[%v]", enumType.ToString()))
	}
	Obj.CallbackPoll[enumType] = cbVal
	Obj.conn.Subscribe(enumType.ToString(), func(msg *nats.Msg) {
		coreMsg := LoadCoreMsgStream(msg.Data)
		theKey := fmt.Sprintf("Msg:%v", coreMsg.Uid)
		success, err := Obj.RedisCli.SetNX(context.Background(), theKey, "", 60*time.Second).Result()
		if !success {
			log.Printf("[SetNX] double spend [%v]", theKey)
			backByteStr, _ := Obj.RedisCli.GetEx(context.Background(), theKey, 60*time.Second).Result()
			msg.Respond([]byte(backByteStr))
			return
		}
		if err != nil {
			log.Printf("[SetNX][%v][%v]", theKey, err.Error())
			return
		}
		refType := cbVal.Type().In(0)
		refVal := reflect.New(refType.Elem())
		jsoniter.Unmarshal(coreMsg.Stream, refVal.Interface())
		outParam := cbVal.Call([]reflect.Value{refVal})
		var backByte []byte
		switch len(outParam) {
		case 0:
			backByte = []byte("OK")
		case 1:
			tmpBackBytes, _ := jsoniter.Marshal(outParam[0].Interface())
			backByte = tmpBackBytes
		}
		setExStr, err := Obj.RedisCli.SetEX(context.Background(), theKey, backByte, 60*time.Second).Result()
		if err != nil {
			log.Printf("[SetNX]set back [%v] err [%v]", theKey, err.Error())
			return
		}
		log.Printf("[SetNX] setExStr[%v]", setExStr)

		msg.Respond(backByte)
	})
}
