//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

package app

import (
	"fmt"
	// "reflect"
	"github.com/adamkdean/consul-network-poc/utils/pkg/consul"
	"github.com/satori/go.uuid"
)

type Instance struct {
	ID     string
	Consul *consul.Instance
}

func (i *Instance) Init(addr string) {
	i.Consul = consul.New()
	i.Must(i.Consul.Initialize(addr))
	i.Must(i.Consul.RegisterService(i.ID, "stargate", "INITIALIZED"))
}

func (i *Instance) Must(err error) {
	if err != nil {
		// Log error? Recover?
		panic(err.Error())
	}
}

func New() *Instance {
	// Generate a UUID using V1 which incorporates both
	// timestamp and MAC address, and convert to string
	uuid := fmt.Sprintf("%s", uuid.Must(uuid.NewV1()))

	return &Instance{
		ID: uuid,
	}
}
