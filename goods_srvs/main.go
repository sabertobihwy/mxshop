package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/goods_srvs/global"
	"mxshop_srvs/goods_srvs/handler"
	"mxshop_srvs/goods_srvs/initialize"
	"mxshop_srvs/goods_srvs/proto"
	"mxshop_srvs/goods_srvs/utils"
)

var serviceId string

func Register(consuladdr string, grpcHost string, port int, name string, tags []string, id string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consuladdr, 8500)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// health check instance
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", grpcHost, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	serviceId = fmt.Sprintf("%s", uuid.NewV4())
	// registration instance
	regis := &api.AgentServiceRegistration{
		Name:    name,
		ID:      serviceId,
		Port:    port,
		Tags:    tags,
		Address: grpcHost,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(regis)
	if err != nil {
		panic(err)
	}
	return client, err

}

func main() {
	// initialize zap
	initialize.InitilizeLogger()
	// initialize config
	initialize.InitConfig()
	// initialize db
	initialize.InitDB()

	IP := flag.String("ip", "0.0.0.0", "ip addr") // default: 0.0.0.0
	PORT := flag.Int("port", 50051, "port")
	flag.Parse()
	if *PORT == 0 {
		*PORT, _ = utils.GetFreePort()
	}
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *PORT)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{}) // grpc-server, implemented server
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	client, err := Register(global.ServiceConfig.ConsulConfig.Host, global.ServiceConfig.Host, *PORT,
		global.ServiceConfig.Name, global.ServiceConfig.Tags, global.ServiceConfig.Name)
	if err != nil {
		zap.S().Debugf(err.Error())
		panic("register error")
	}
	fmt.Println("=========" + fmt.Sprintf("%s:%d", *IP, *PORT))
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
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Debugf(fmt.Sprintf("deregister error"))
		panic(err)
	}
}
