package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/user_srvs/global"
	"mxshop_srvs/user_srvs/handler"
	"mxshop_srvs/user_srvs/initialize"
	"mxshop_srvs/user_srvs/proto"
)

func Register(consuladdr string, grpcHost string, port int, name string, tags []string, id string) error {
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
	// registration instance
	regis := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: grpcHost,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(regis)
	if err != nil {
		panic(err)
	}
	return err

}

func main() {
	// initialize zap
	initialize.InitilizeLogger()
	// initialize config
	initialize.InitConfig()
	// initialize db
	initialize.InitDB()

	IP := flag.String("ip", global.ServiceConfig.Host, "ip addr") // default: 0.0.0.0
	PORT := flag.Int("port", int(global.ServiceConfig.ConsulConfig.Port), "port")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *PORT)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{}) // grpc-server, implemented server
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	err := Register(global.ServiceConfig.ConsulConfig.Host, global.ServiceConfig.Host, global.ServiceConfig.ConsulConfig.Port,
		global.ServiceConfig.Name, []string{"mxshop", "bobby"}, global.ServiceConfig.Name)
	if err != nil {
		zap.S().Debugf(err.Error())
		panic("register error")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *PORT))
	if err != nil {
		zap.S().Debugf(fmt.Sprintf("%s:%d", *IP, *PORT))
		panic("listen error")
	}
	err = server.Serve(lis)
	if err != nil {
		panic("fail to start grpc")
	}
}
