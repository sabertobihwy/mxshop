package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/order_srvs/global"
	"mxshop_srvs/order_srvs/handler"
	"mxshop_srvs/order_srvs/initialize"
	"mxshop_srvs/order_srvs/proto"
	"mxshop_srvs/order_srvs/utils"
	"mxshop_srvs/order_srvs/utils/register/consul"
)

func main() {
	// initialize zap
	initialize.InitilizeLogger()
	// initialize config
	initialize.InitConfig()
	// init srvs
	initialize.InitSrvCli()
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
	proto.RegisterOrderServer(server, &handler.OrderServer{}) // grpc-server, implemented server
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
	// ordersrv listen to "order_timeout" mq
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("order"), // sequentially consume message; load balance
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.2.112:9876"})),
	)
	err = c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.AutoTimeout)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	//quit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = client.Agent().ServiceDeregister(consul.ServiceId); err != nil {
		zap.S().Debugf(fmt.Sprintf("deregister error"))
		panic(err)
	}
}
