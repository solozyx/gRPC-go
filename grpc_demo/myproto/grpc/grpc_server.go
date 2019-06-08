package main

import (
	"context"
	"fmt"
	"net"

	pb "gRPC-go/grpc_demo/myproto"
	"google.golang.org/grpc"
)

// 定义1个服务对象
type server struct{}

// rpc
// 函数关键字 (对象) 函数名 (客户端发送过来的参数,返回给客户端的内容) 错误返回值

// grpc
// 函数关键字 (对象) 函数名 (context,客户端发送过来的参数) (返回给客户端的内容,错误返回值)

// 接口方法必须全部实现
// type HelloServerServer interface {
// SayHello(context.Context, *HelloReq) (*HelloRsp, error)
// SayName(context.Context, *NameReq) (*NameRsp, error)
// }

func (s *server) SayHello(ctx context.Context, in *pb.HelloReq) (out *pb.HelloRsp, err error) {
	return &pb.HelloRsp{Msg: "grpc_server SayHello " + in.Name}, nil
}

func (s *server) SayName(ctx context.Context, in *pb.NameReq) (out *pb.NameRsp, err error) {
	return &pb.NameRsp{Msg: "grpc_server SayName " + in.Name}, nil
}

func main() {
	// 创建网络
	ln, err := net.Listen("tcp", ":10086")
	if err != nil {
		fmt.Println("网络错误", err)
	}
	defer ln.Close()

	// 创建grpc的服务
	srv := grpc.NewServer()

	// 注册grpc服务
	pb.RegisterHelloServerServer(srv, &server{})

	// 创建网络连接监听
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("网络错误", err)
	}
}
