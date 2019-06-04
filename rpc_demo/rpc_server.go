package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

/**
- 方法是导出的
- 方法有两个参数，都是导出类型或内建类型
- 方法的第二个参数是指针
- 方法只有一个error接口类型的返回值
func (t *T) MethodName(argType T1, replyType *T2) error
*/

// 注册1个对象
type Panda int

// 对象注册方法
// 函数关键字func (对象) 函数名方法名 (对端发送数据,本端返回数据) 错误接口
func (p *Panda) GetInfo(argType int, replyType *int) error {
	fmt.Println("rpc client 发送的数据 = ", argType)
	*replyType = argType + 12306
	return nil
}

func main() {
	// 创建一个对象
	p := new(Panda)
	// 服务端注册对象
	rpc.Register(p)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":10010")
	if err != nil {
		log.Fatalln("网络错误")
	}
	defer l.Close()
	http.Serve(l, nil)
}
