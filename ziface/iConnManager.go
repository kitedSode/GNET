package ziface

type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connId uint32) (IConnection, error) // 利用ConnId获取连接
	Len() int                               // 获取当前
	ClearConn()                             // 删除并停止所有连接
}
