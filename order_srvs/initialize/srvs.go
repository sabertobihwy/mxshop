package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"log"

	"mxshop_srvs/order_srvs/global"
	"mxshop_srvs/order_srvs/proto"
)

func InitSrvCli() {
	var err error
	// goods cli
	conn, err := grpc.Dial(fmt.Sprintf("consul://%s:8500/%s?wait=14s&tag=%s", global.ServiceConfig.ConsulConfig.Host, global.ServiceConfig.GoodSrv.Name, "mxshop"),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.GoodsClient = proto.NewGoodsClient(conn)
	// inventory cli
	conn, err = grpc.Dial(fmt.Sprintf("consul://%s:8500/%s?wait=14s&tag=%s", global.ServiceConfig.ConsulConfig.Host, global.ServiceConfig.InventorySrv.Name, "mxshop"),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	global.StocksClient = proto.NewStocksClient(conn)
}
