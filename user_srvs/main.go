package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/user_srvs/handler"
	"mxshop_srvs/user_srvs/proto"
	"net"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip addr") // default: 0.0.0.0
	PORT := flag.String("port", "50051", "port")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *PORT)
	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handler.UserServer{}) // grpc-server, implemented server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *IP, *PORT))
	if err != nil {
		panic("listen error")
	}
	err = g.Serve(lis)
	if err != nil {
		panic("fail to start grpc")
	}
}
