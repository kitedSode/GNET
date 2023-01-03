package znet

type Message struct {
	Id      uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// NewMsgPackage 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgId 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetMsgLen 设置消息的长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

// SetMsgData 设置消息的内容
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
