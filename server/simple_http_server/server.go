package main

import (
	"context"
	"gRPC-go/pkg/gtls"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"

	pb "gRPC-go/proto"
)

const PORT = "9003"

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
}

func GetHTTPServeMux() *http.ServeMux {
	// http.NewServeMux：创建一个新的 ServeMux，ServeMux 本质上是一个路由表
	// 它默认实现了 ServeHTTP，因此返回 Handler 后
	// 可直接通过 HandleFunc 注册 pattern 和处理逻辑的方法
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("grpc_http_server: go-grpc-example"))
	})
	return mux
}

func main() {
	certFile := "../../conf/server/server.pem"
	keyFile := "../../conf/server/server.key"
	tlsServer := gtls.Server{
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	c, err := tlsServer.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsServer.GetTLSCredentials err: %v", err)
	}

	// HTTP/1.1 mux
	mux := GetHTTPServeMux()

	// HTTP/2 gRPC server
	server := grpc.NewServer(grpc.Creds(c))

	pb.RegisterSearchServiceServer(server, &SearchService{})

	// http.ListenAndServeTLS：可简单的理解为提供监听 HTTPS 服务的方法
	// 协议判断转发 判断 -> 转发 -> 响应
	http.ListenAndServeTLS(":"+PORT,
		certFile,
		keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}),
	)

	//http.ListenAndServe(":"+PORT,
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
	//			server.ServeHTTP(w, r)
	//		} else {
	//			mux.ServeHTTP(w, r)
	//		}
	//	}),
	//)
}
