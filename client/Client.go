package main

import (
	znet2 "GNET/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(time.Second)

	// 1.直接连接到远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}
	defer conn.Close()

	for {
		// 发送封包数据
		dp := znet2.NewDataPack()
		msg := znet2.NewMsgPackage(0, []byte("ZinxV0.8 client msg"))
		if binaryData, err := dp.Pack(msg); err != nil {
			fmt.Println("pack err:", err)
			return
		} else {
			conn.Write(binaryData)
		}

		// 服务器应该给我们回复一个message数据，MsgId:1 doing ping

		// 先读取流中的head部分来得到Id和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(conn, binaryHead)
		if err != nil {
			panic("read binaryHead error:" + err.Error())
		}
		callBackHead, err := dp.Unpack(binaryHead)
		if err != nil {
			panic("unpack headData error:" + err.Error())
		}

		// 再根据DataLen进行第二次读取，将data读出来
		if callBackHead.GetMsgLen() > 0 {
			receiveMsg := callBackHead.(*znet2.Message)
			receiveMsg.SetMsgData(make([]byte, callBackHead.GetMsgLen()))

			_, err = io.ReadFull(conn, receiveMsg.GetData())
			if err != nil {
				panic("read back data error:" + err.Error())
			}

			fmt.Println("----> Receive Server Msg : msgId =", receiveMsg.GetMsgId(),
				"msgLen =", receiveMsg.GetMsgLen(), "data =", string(receiveMsg.GetData()))
		}
		time.Sleep(time.Second)
	}

}
