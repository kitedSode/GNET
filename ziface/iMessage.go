package ziface

// IMessage 将请求的消息封装到一个Message中，定义抽象的接口
type IMessage interface {
	// GetMsgId 获取消息的ID
	GetMsgId() uint32

	// GetMsgLen 获取消息的长度
	GetMsgLen() uint32

	// GetData 获取消息的内容
	GetData() []byte

	// SetMsgId 设置消息的ID
	SetMsgId(uint32)

	// SetMsgLen 设置消息的长度
	SetMsgLen(uint32)

	// SetMsgData 设置消息的内容
	SetMsgData([]byte)
}
