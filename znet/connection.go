package znet

import (
	ziface2 "GNET/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TcpServer  ziface2.IServer     // 当前Conn属于哪个Server，在Conn初始化化时添加
	Conn       *net.TCPConn        // 当前连接的socket TCP套接字
	ConnID     uint32              // 连接的ID
	isClosed   bool                // 当前的连接状态
	ExitChan   chan bool           // 告知当前连接已经退出或停止的 channel
	MsgHandler ziface2.IMsgHandler // 该连接处理的方法Router
	MsgChan    chan []byte         // 无缓冲的管道，用于读、写Goroutine之间的消息通信

	property     map[string]interface{} // 连接属性
	propertyLock sync.RWMutex           // 保护连接属性修改的锁
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface2.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface2.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		MsgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer func() {
		fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.Conn.RemoteAddr())
		c.Stop()
	}()
	// 次数关闭
	//var times = 0

	// 创建一个拆包解包对象
	dp := NewDataPack()
	for {
		// 读取客户端的Msg Head 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:", err)
			break
		}
		//fmt.Println("=========> times = ", times)
		//times ++
		//if times == 3{
		//	break
		//}
		// 拆包，得到msgId 和 msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		// 根据dataLen再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("get msg data err:", err)
				break
			}
		}
		msg.SetMsgData(data)

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 将消息交由Worker工作池来处理
		c.MsgHandler.SendMsgToTaskQueue(req)
	}
}

// SendMsg 提供一个SendMsg方法将我们要发送的消息
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when sending msg")
	}
	// 获得封包对象
	dp := NewDataPack()

	// 封包后的消息格式 MsgDataLen/MsgId/Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error, msgId =", msgId)
		return errors.New("pack error msg")
	}

	c.MsgChan <- binaryMsg
	return nil
}

// StartWriter 写消息Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.MsgChan:
			// 发送数据给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader以及退出，此时Writer也要退出
			return
		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动从当前连接的读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 用户自定义的连接建立时的业务，执行钩子函数
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 执行用户自定义的连接断开时的钩子函数
	c.TcpServer.CallOnConnStop(c)

	// close socket connect
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	if _, err := c.TcpServer.GetConnMgr().Get(c.ConnID); err == nil {
		// 将连接从连接管理器中删除
		c.TcpServer.GetConnMgr().Remove(c)
	}

	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)
}

// GetTCPConnection 获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SetProperty 创建连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; !ok {
		return nil, errors.New("key is not exist")
	} else {
		return value, nil
	}

}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
