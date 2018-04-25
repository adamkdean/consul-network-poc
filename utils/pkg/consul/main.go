//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

package consul

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

type Instance struct {
	Agent   *api.Agent
	Catalog *api.Catalog
	Client  *api.Client
	Config  *api.Config
	KV      *api.KV
}

func (i *Instance) Initialize(addr string) error {
	i.Config = api.DefaultConfig()
	i.Config.Address = addr
	client, err := api.NewClient(i.Config)
	if err != nil {
		return err
	}

	i.Client = client
	i.Agent = i.Client.Agent()
	i.Catalog = i.Client.Catalog()
	i.KV = i.Client.KV()
	return nil
}

func (i *Instance) RegisterService(id, name, tag string) error {
	s := &api.AgentServiceRegistration{
		ID:   id,
		Name: name,
		Tags: []string{tag},
	}
	return i.Agent.ServiceRegister(s)
}

func (i *Instance) DeregisterService(id string) error {
	return i.Agent.ServiceDeregister(id)
}

func (i *Instance) GetAllServices() (map[string][]string, error) {
	s, _, err := i.Catalog.Services(nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (i *Instance) GetService(service, tag string) ([]*api.CatalogService, error) {
	s, _, err := i.Catalog.Service(service, tag, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (i *Instance) KeyExists(key string) bool {
	kvp, _, err := i.KV.Get(key, nil)
	if err != nil || kvp == nil {
		return false
	}

	return true
}

func (i *Instance) WriteStringToKey(key, data string) error {
	kvp := &api.KVPair{Key: key, Value: []byte(data)}
	_, err := i.KV.Put(kvp, nil)
	return err
}

func (i *Instance) WriteStructToKey(key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	kvp := &api.KVPair{Key: key, Value: b}
	_, err = i.KV.Put(kvp, nil)
	return err
}

func New() *Instance {
	return &Instance{}
}
