package ziface

/**
消息管理抽象层
*/

type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // DoMsgHandler 调度或执行对应的Router消息处理方法
	AddRouter(msgId uint32, router IRouter) // AddRouter 为消息分配具体的处理逻辑
	StartWorkerPool()                       // 启动Worker工作池
	SendMsgToTaskQueue(request IRequest)    // 将消息交给TaskQueue, 由Worker进行处理
}
