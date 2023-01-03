package znet

import (
	"GNET/utils"
	ziface2 "GNET/ziface"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name       string               //服务器的名称
	IPVersion  string               //服务器绑定的ip版本
	IP         string               //服务器监听的IP
	Port       int                  // 服务器监听的端口
	MsgHandler ziface2.IMsgHandler  //当前的Server的消息管理模块，用来绑定MsgId和对应的处理业务API关系
	ConnMgr    ziface2.IConnManager // 当前Server的连接管理器

	OnConnStart func(conn ziface2.IConnection) // 该Server的连接创建时的Hook函数
	OnConnStop  func(conn ziface2.IConnection) // 该Server的连接断开时的Hook函数
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port : %d is starting\n",
		s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version %s, MaxConn : %d, MaxPacketSize : %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		errCheck("resolve tcp addr", err)

		// 2.监听服务器的地址
		listen, err := net.ListenTCP(s.IPVersion, addr)
		errCheck("Listen", err)
		defer listen.Close()
		fmt.Println("start Zinx server success,", s.Name, "success, Listening...")

		var cid uint32
		cid = 0
		// 3.阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listen.AcceptTCP()
			// TODO 我认为这里需要再单独启动一个Goroutine来负责处理连接业务
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 判断当前的连接个数是否超过最大连接限制，如果超过则关闭新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				fmt.Println("too many connections!")
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go dealConn.Start()

		}
	}()

}

func (s *Server) Stop() {
	// 释放服务器的资源，状态或者一些已经开辟的连接信息
	fmt.Println("[STOP] Zinx server name", s.Name)
	//TODO 执行此方法会引起死锁！（ConnManager.ClearConn(Lock) -> Connection.Stop -> ConnManager.Remove(Lock)）
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 其他业务

	//阻塞状态，接收信号会结束程序
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("The server is down")
}

// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface2.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Success!!")
}

func (s *Server) GetConnMgr() ziface2.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface2.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop SetOnConnStart 设置该Server的连接断开时Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface2.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook 函数
func (s *Server) CallOnConnStart(conn ziface2.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-----> CallOnConnStart")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook 函数
func (s *Server) CallOnConnStop(conn ziface2.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> CallOnConnStop")
		s.OnConnStop(conn)
	}
}

// NewServer 初始化Server模块的方法
func NewServer() ziface2.IServer {
	global := utils.GlobalObject
	s := &Server{
		Name:       global.Name,
		IPVersion:  "tcp4",
		IP:         global.Host,
		Port:       global.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

func errCheck(info string, err error) {
	if err != nil {
		fmt.Println("info:", info, "err:", err)
		os.Exit(-1)
	}
}
