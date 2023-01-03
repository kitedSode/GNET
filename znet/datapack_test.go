package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是测试datapack拆包 封包的单元测试
func TestDataPack_Pack(t *testing.T) {
	// 1.创建SocketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// 创建一个go 来负责从客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println("ok了家人们")
			go func(client net.Conn) {
				defer client.Close()
				// 处理客户端的请求
				// ------> 拆包的过程 <------
				// 定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(client, headData)
					if err != nil {
						fmt.Println("read head error:", err)
						return
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("Unpack error:", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// msg是有数据的，需要进行第二次读取
						// 第二次从conn读，根据head中的dataLen读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						// 完整的一个消息已经读取完毕
						fmt.Printf("----> Receive MsgId :%d, dataLen = %d, data = %s\n", msg.Id, msg.DataLen, string(msg.GetData()))
					}

				}
			}(conn)
		}
	}()
	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7000")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendDatal1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("msg1 pack error:", err)
		return
	}

	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'0', '1', '2', '3', '4', '5', '6'},
	}
	sendDatal2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("msg1 pack error:", err)
		return
	}

	// 将两个包粘在一起
	sendDatal1 = append(sendDatal1, sendDatal2...)

	// 一次性发送给服务端
	conn.Write(sendDatal1)

	select {}
}
