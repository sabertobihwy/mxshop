package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/stocks_srvs/global"
	"mxshop_srvs/stocks_srvs/handler"
	"mxshop_srvs/stocks_srvs/initialize"
	"mxshop_srvs/stocks_srvs/proto"
	"mxshop_srvs/stocks_srvs/utils"
	"mxshop_srvs/stocks_srvs/utils/register/consul"
)

func main() {
	// initialize zap
	initialize.InitilizeLogger()
	// initialize config
	initialize.InitConfig()
	// initialize db
	initialize.InitDB()
	initialize.InitializeRedis(global.ServiceConfig.RedisConfig.Host, global.ServiceConfig.RedisConfig.Port)

	IP := flag.String("ip", global.ServiceConfig.Host, "ip addr") // default: 0.0.0.0
	PORT := flag.Int("port", 50051, "port")
	flag.Parse()
	if *PORT == 0 {
		*PORT, _ = utils.GetFreePort()
	}
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *PORT)

	server := grpc.NewServer()
	proto.RegisterStocksServer(server, &handler.StocksServer{}) // grpc-server, implemented server
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	client, err := consul.Register(global.ServiceConfig.ConsulConfig.Host, global.ServiceConfig.Host, *PORT,
		global.ServiceConfig.Name, global.ServiceConfig.Tags, global.ServiceConfig.Name)
	if err != nil {
		zap.S().Debugf(err.Error())
		panic("register error")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		zap.S().Debugf(fmt.Sprintf("%s:%d", *IP, *PORT))
		panic("listen error")
	}
	go func() {
		err = server.Serve(lis) // block -> turn into goroutine
		if err != nil {
			panic("fail to start grpc")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(consul.ServiceId); err != nil {
		zap.S().Debugf(fmt.Sprintf("deregister error"))
		panic(err)
	}
}
