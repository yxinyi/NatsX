package natsx

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

type TestMsg struct {
	Msg string
}

type TestMsgBack struct {
	Msg string
}

type CoreMsg struct {
	Uid    int64
	Stream []byte
}

var uuidCount = int64(0)

func getUid() int64 {
	uuidCount++
	return time.Now().Unix()<<32 + uuidCount
}
func ToCoreMsgStream(msg interface{}) []byte {
	stream, _ := jsoniter.Marshal(msg)
	cm := &CoreMsg{
		Uid:    getUid(),
		Stream: stream,
	}
	cmStream, _ := jsoniter.Marshal(cm)
	return cmStream
}

func LoadCoreMsgStream(cmStream []byte) *CoreMsg {
	cm := &CoreMsg{}
	jsoniter.Unmarshal(cmStream, cm)
	return cm
}
