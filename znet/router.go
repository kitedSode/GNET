package znet

import (
	"GNET/ziface"
)

type BaseRouter struct {
}

/**
这里的BaseRouter的方法都为空是因为有的Router不希望有PreHandle、PostHandle这两个业务。
所以Router全部继承BaseRouter的好处就是不需要实现PreHandle和PostHandle
*/

// PreHandle 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

// Handle 在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
