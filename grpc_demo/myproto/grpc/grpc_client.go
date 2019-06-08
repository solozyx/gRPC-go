package main

import (
	"context"
	"fmt"

	pb "gRPC-go/grpc_demo/myproto"
	"google.golang.org/grpc"
)

// type HelloServerClient interface {
// SayHello(ctx context.Context, in *HelloReq, opts ...grpc.CallOption) (*HelloRsp, error)
// SayName(ctx context.Context, in *NameReq, opts ...grpc.CallOption) (*NameRsp, error)
// }

func main() {
	// grpc客户端 连接 grpc服务器
	conn, err := grpc.Dial("127.0.0.1:10086", grpc.WithInsecure())
	if err != nil {
		fmt.Println("网络异常", err)
	}
	// 网络延迟关闭
	defer conn.Close()

	// 获得grpc句柄
	c := pb.NewHelloServerClient(conn)

	// 通过句柄调用函数
	helloRsp, err := c.SayHello(context.Background(), &pb.HelloReq{Name: "熊猫"})
	if err != nil {
		fmt.Println("grpc_client 调用 grpc_server SayHello 服务失败")
	}
	fmt.Println("grpc_client 调用 grpc_server SayHello 服务的返回", helloRsp.Msg)

	nameRsp, err := c.SayName(context.Background(), &pb.NameReq{Name: "托尼斯塔克"})
	if err != nil {
		fmt.Println("grpc_client 调用 grpc_server SayName 服务失败")
	}
	fmt.Println("grpc_client 调用 grpc_server SayName 服务的返回", nameRsp.Msg)
}
