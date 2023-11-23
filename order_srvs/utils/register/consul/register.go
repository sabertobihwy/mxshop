package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
)

var ServiceId string

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
	ServiceId = fmt.Sprintf("%s", uuid.NewV4())
	// registration instance
	regis := &api.AgentServiceRegistration{
		Name:    name,
		ID:      ServiceId,
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
