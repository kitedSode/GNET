package znet

import (
	ziface2 "GNET/ziface"
)

type Request struct {
	conn ziface2.IConnection // 已经和客户端建立好的连接
	msg  ziface2.IMessage    // 客户端请求的数据
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() ziface2.IConnection {
	return r.conn
}

// GetData 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
