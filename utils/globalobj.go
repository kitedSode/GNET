package utils

import (
	"GNET/ziface"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
)

/**
存储一切有关Zinx框架有关的全局参数，供其他模块使用
一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	/**
	Server
	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的的端口号
	Name      string         // 当前服务器的名称

	/**
	Zinx
	*/
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 业务工作Worker池的数量
	MaxWorkerTaskLen uint32 // 业务工作Worker对应负责的任务队列最大任务存储数量

	/*
		config file path
	*/
	ConfFilePath string // 配置文件地址

}

// GlobalObject 定义一个全局对外的GlobalObj
var GlobalObject *GlobalObj

// Reload 从zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	fmt.Println(g.ConfFilePath)
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	// TODO
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject
func init() {

	// 如果配置文件没有加载则使用该配置
	GlobalObject = &GlobalObj{
		Name:             "GNET",
		Version:          "V0.10",
		TcpPort:          8080,
		Host:             "127.0.0.1",
		MaxConn:          10,
		MaxPackageSize:   4096,
		ConfFilePath:     "./conf/setting.json",
		WorkerPoolSize:   uint32(runtime.NumCPU()),
		MaxWorkerTaskLen: 1024,
	}

	// 获取配置文件中用户自定义的参数
	GlobalObject.Reload()
}
