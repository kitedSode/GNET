package main

import (
	ziface2 "GNET/ziface"
	znet2 "GNET/znet"
	"fmt"
)

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet2.BaseRouter
}

// Handle Test Handle
func (pr *PingRouter) Handle(request ziface2.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Println("receive from client: msgId =", request.GetMsgId(),
		", data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("doing ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet2.BaseRouter
}

// Handle for HelloRouter
func (hr *HelloRouter) Handle(request ziface2.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Println("receive from client: msgId =", request.GetMsgId(),
		", data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello...."))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(connection ziface2.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called...")
	if err := connection.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	connection.SetProperty("Name", "Flacus")
	connection.SetProperty("QQ", "1054506646")
	connection.SetProperty("Phone", "15856385519")
}

// DoConnectionLost 创建连接之后执行的钩子函数
func DoConnectionLost(connection ziface2.IConnection) {

	fmt.Println("====> DoConnectionLost is Called...")
	//if err := connection.SendMsg(202,[]byte("DoConnection LOST")); err != nil{
	//	fmt.Println(err)
	//}

	if value, err := connection.GetProperty("Name"); err == nil {
		fmt.Println("Name:", value)
	}
	if value, err := connection.GetProperty("QQ"); err == nil {
		fmt.Println("QQ:", value)
	}
	if value, err := connection.GetProperty("Phone"); err == nil {
		fmt.Println("Phone:", value)
	}

}

func main() {
	//1 创建一个server句柄，使用Zinx的api
	s := znet2.NewServer()

	//2 给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//3 设置 Server的钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//4 启动server
	s.Serve()
}
