package main

import (
	"fmt"
	"log"
	"natsx/natsx"
)

func main() {
	natsx.Obj.Connection("nats://127.0.0.1:4222")
	//natsx.Obj.StartLoop()
	natsx.RegisterMsg(natsx.EnumTypeTest, func(msg *natsx.TestMsg) {
		log.Printf("[%v][%v]", natsx.EnumTypeTest.ToString(), msg.Msg)
	})

	for idx := 0; idx < 100000; idx++ {
		//重复发送
		msg := &natsx.TestMsg{
			Msg: fmt.Sprintf("%v", idx),
		}
		natsx.SyncMsg(natsx.EnumTypeTest, msg, nil)
	}

}
