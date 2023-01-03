package znet

import (
	"GNET/ziface"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                  // 读写连接的读写锁
}

// NewConnManager 创建一个连接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源Map 加写锁(其他go程无法加读锁或写锁)
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully;totally conn's num =", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源Map 加读锁(其他go程无法加写锁)
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	// 删除连接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnId =", conn.GetConnID(), "successfully;totally conn's num =", cm.Len())
}

func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源Map 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connId = [" + strconv.Itoa(int(connId)) + "] is not exist")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	// 保护共享资源Map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 停止并删除全部的连接信息
	for connId, conn := range cm.connections {
		// 停止
		conn.Stop()

		// 删除
		delete(cm.connections, connId)
	}

	fmt.Println("Clear All Connections successfully; conn num =", cm.Len())
}
