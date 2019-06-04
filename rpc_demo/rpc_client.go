package main

import (
	"log"
	"net/rpc"
)

func main() {
	// 创建网络连接
	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:10010")
	if err != nil {
		log.Println("网络连接失败")
	}
	defer cli.Close()

	var output int
	// 调用远程方法
	err = cli.Call("Panda.GetInfo", 10086, &output)
	if err != nil {
		log.Println("远程调用失败")
	}

	log.Println(output)
}
