package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端连接处理使用
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 得到连接管理
	GetConnMgr() IConnManager
	// SetOnConnStart 设置该Server的连接创建时Hook函数
	SetOnConnStart(fun func(IConnection))
	// SetOnConnStop SetOnConnStart 设置该Server的连接断开时Hook函数
	SetOnConnStop(fun func(IConnection))
	// CallOnConnStart 调用连接OnConnStart Hook 函数
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用连接OnConnStop Hook 函数
	CallOnConnStop(connection IConnection)
}
