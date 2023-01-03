package znet

import (
	"GNET/utils"
	ziface2 "GNET/ziface"
	"fmt"
	"log"
)

/**
消息处理模块的具体实现
*/

type MsgHandler struct {
	Apis           map[uint32]ziface2.IRouter // 存放每个MsgId所对应的处理方法
	WorkerPoolSize uint32                     // 业务工作Worker池的数量
	TaskQueue      []chan ziface2.IRequest    // Worker负责取任务的消息队列
	arr            []int
}

// NewMsgHandler 初始化/创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface2.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface2.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度或执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface2.IRequest) {
	// 从Request的msgId来得到相应的router
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId =", request.GetMsgId(), "is NOT FOUND ! Need Register !")
		return
	}

	// 根据MsgId来调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息分配具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface2.IRouter) {
	// 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		log.Fatalf("repeat api, msgId = %d\n", msgId)
	}

	// 添加msg与API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("add api MsgId = ", msgId, "success !")
}

// StartWorkerPool 启动Worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	// TODO 新增修改（之前没有此步骤）
	mh.TaskQueue = make([]chan ziface2.IRequest, mh.WorkerPoolSize)
	// 遍历需要启动Worker的数量，依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个Worker被启动
		// 给当前Worker对用的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface2.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

// StartOneWorker 启动一个Worker工作流程
func (mh MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface2.IRequest) {
	fmt.Println("WorkerId = ", workerId, "is started")

	// 不断地等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}

}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface2.IRequest) {
	// 根据ConnId来分配当前的连接应该由哪个Worker负责处理
	// 轮询的平均分配法则

	// 得到需要处理此条连接的WorkerId
	connId := request.GetConnection().GetConnID()
	workerId := connId % mh.WorkerPoolSize
	fmt.Println("Add ConnId =", connId,
		"request msgId =", request.GetMsgId(),
		"to WorkerId =", workerId)
	// 将请求消息发送给任务队列
	mh.TaskQueue[workerId] <- request
}
