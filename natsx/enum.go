package natsx

type Enum string

func (e Enum) ToString() string {
	return string(e)
}

const (
	MsgTypeSync    = 1
	MsgTypeASync   = 2
	MsgTypePublish = 3
)
const (
	EnumTypeTest Enum = "Test"
)
