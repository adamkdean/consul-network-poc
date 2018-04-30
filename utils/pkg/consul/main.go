//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

// Package consul provides helper methods to work with Consul
// in context of the DADI Cloud decentralized network.
package consul

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
)

// Instance ...
type Instance struct {
	Agent   *api.Agent
	Catalog *api.Catalog
	Client  *api.Client
	Config  *api.Config
	KV      *api.KV
}

// Initialize ...
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

// RegisterService ...
func (i *Instance) RegisterService(id, name, tag string) error {
	s := &api.AgentServiceRegistration{
		ID:   id,
		Name: name,
		Tags: []string{tag},
	}
	return i.Agent.ServiceRegister(s)
}

// DeregisterService ...
func (i *Instance) DeregisterService(id string) error {
	return i.Agent.ServiceDeregister(id)
}

// GetAllServices ...
func (i *Instance) GetAllServices() (map[string][]string, error) {
	s, _, err := i.Catalog.Services(nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetService ...
func (i *Instance) GetService(service, tag string) ([]*api.CatalogService, error) {
	s, _, err := i.Catalog.Service(service, tag, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetServiceManifest ...
func (i *Instance) GetServiceManifest(sv, id string) (*ServiceManifest, error) {
	// Get a list of key value pairs for service prefix
	key := fmt.Sprintf("%s/%s", sv, id)
	kvp, _, err := i.KV.Get(key, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("kvp: %v\n", kvp)

	// Parse the key value pair into a ServiceManifest struct
	m := &ServiceManifest{}
	json.Unmarshal(kvp.Value, &m)
	return m, nil
}

// GetServiceManifests ...
func (i *Instance) GetServiceManifests(sv string) ([]*ServiceManifest, error) {
	// Get a list of key value pairs for service prefix
	prefix := fmt.Sprintf("%s/", sv)
	kvp, _, err := i.KV.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	// Parse the key value pairs into ServiceManifest structs
	manifests := []*ServiceManifest{}
	for i := range kvp {
		m := &ServiceManifest{}
		json.Unmarshal(kvp[i].Value, &m)
		manifests = append(manifests, m)
	}

	return manifests, nil
}

// RemoveServiceManifest ...
func (i *Instance) RemoveServiceManifest(service, id string) error {
	key := fmt.Sprintf("%s/%s", service, id)
	_, err := i.KV.Delete(key, nil)
	return err
}

// KeyExists ...
func (i *Instance) KeyExists(key string) bool {
	kvp, _, err := i.KV.Get(key, nil)
	if err != nil || kvp == nil {
		return false
	}
	return true
}

// WriteStringToKey ...
func (i *Instance) WriteStringToKey(key, data string) error {
	kvp := &api.KVPair{Key: key, Value: []byte(data)}
	_, err := i.KV.Put(kvp, nil)
	return err
}

// WriteStructToKey ...
func (i *Instance) WriteStructToKey(key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	kvp := &api.KVPair{Key: key, Value: b}
	_, err = i.KV.Put(kvp, nil)
	return err
}

// New ...
func New() *Instance {
	return &Instance{}
}
