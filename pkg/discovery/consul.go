package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"strconv"
	"strings"
)

func InitServiceDiscovery() {
	log.Print("initializing consul client")

	consulClient, err := api.NewClient(
		&api.Config{
			Address: os.Getenv("CONSUL_HOST"),
		},
	)
	if err != nil {
		log.Fatal("error creating consul client")
	}

	log.Print("register service in consul")
	agent := consulClient.Agent()
	parsedPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("port parse error")
	}

	tags := strings.Split(os.Getenv("CONSUL_TAGS"), ",")
	service := &api.AgentServiceRegistration{
		Name:    os.Getenv("CONSUL_SERVICE_NAME"),
		Port:    parsedPort,
		Address: os.Getenv("CONSUL_SERVICE_ADDRESS"),
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", os.Getenv("CONSUL_SERVICE_ADDRESS"), parsedPort),
			Interval:                       "10s",
			Timeout:                        "2s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	err = agent.ServiceRegister(service)
	if err != nil {
		log.Fatalf("error while service registration due to error '%s'", err)
	}
	log.Print("service registered in consul")
}
