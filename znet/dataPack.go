package znet

import (
	"GNET/utils"
	"GNET/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// DataPack 封包、拆包的具体模块
type DataPack struct {
}

// NewDataPack 封包拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（字节）+ ID uint32(4字节)
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := &bytes.Buffer{}

	// 将dataLen写进dataBuff中, littleEndian(小端法)
	err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 将MsgId 写进dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 将data数据写进dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到dataLen和MsgId
	msg := Message{}

	// 读dataLen，根据dataLen所定义的类型(uint32)所得它为4个字节的长度，所以这里会读取前四个字节的值
	if err := binary.Read(dataBuff, binary.LittleEndian, msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgId，MsgId所定义的类型也是uint32所以它的长度也是4个字节
	if err := binary.Read(dataBuff, binary.LittleEndian, msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否以及超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return nil, errors.New("too Large msg data receive")
	}

	return &msg, nil
}
