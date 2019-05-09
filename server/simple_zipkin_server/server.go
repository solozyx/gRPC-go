package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"gRPC-go/pkg/gtls"
	pb "gRPC-go/proto"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " gRPC Search Server"}, nil
}

const (
	PORT = "9005"

	SERVICE_NAME              = "simple_zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://192.168.174.131:9411/api/v1/spans"
	ZIPKIN_RECORDER_HOST_PORT = "192.168.174.131:9000"
)

// 初始化 Zipkin 其包含收集器 记录器 跟踪器
// 再利用拦截器在 Server 端实现 SpanContext Payload 的双向读取和管理
func main() {
	// 创建一个 Zipkin HTTP 后端收集器
	collector, err := zipkin.NewHTTPCollector(ZIPKIN_HTTP_ENDPOINT)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err: %v", err)
	}
	// 创建一个基于 Zipkin 收集器的记录器
	recorder := zipkin.NewRecorder(collector, true, ZIPKIN_RECORDER_HOST_PORT, SERVICE_NAME)
	// 创建一个 OpenTracing 跟踪器(兼容 Zipkin Tracer)
	tracer, err := zipkin.NewTracer(
		recorder, zipkin.ClientServerSameSpan(false),
	)
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}

	tlsServer := gtls.Server{
		CaFile:   "../../conf/ca.pem",
		CertFile: "../../conf/server/server.pem",
		KeyFile:  "../../conf/server/server.key",
	}
	c, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			// 返回 grpc.UnaryServerInterceptor 一元拦截器
			// 不同点在于该拦截器会在 gRPC Metadata 中
			// 查找 OpenTracing SpanContext
			// 如果找到则为该服务的 Span Context 的子节点
			//
			// otgrpc.LogPayloads 设置并返回 Option 作用是让 OpenTracing 在双向方向上
			// 记录应用程序的有效载荷payload
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}
