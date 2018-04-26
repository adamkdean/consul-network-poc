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
	"github.com/adamkdean/consul-network-poc/utils/pkg/consul"
	"github.com/adamkdean/consul-network-poc/utils/pkg/state"
	"github.com/satori/go.uuid"
	"time"
)

// Gateway is the application gateway layer for
// the DADI Cloud decentralized network
type Gateway struct {
	ID           string
	Consul       *consul.Instance
	UpdatePeriod int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (g *Gateway) Initialize(addr string) {
	g.Consul = consul.New()
	g.Must(g.Consul.Initialize(addr))
	g.Must(g.UpdateService(state.Initialized))
	go g.UpdateManifest()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (g *Gateway) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := g.Consul.RegisterService(g.ID, "gateway", state); err != nil {
			fmt.Printf("Error registering service: %v (delay %v)\n", err, delay)
			if attempt > maxRetries {
				return fmt.Errorf("Could not register service")
			}

			time.Sleep(time.Duration(delay) * time.Second)
			delay *= 2
		} else {
			fmt.Printf("Successfully registered service with ID %s and state %s\n", g.ID, state)
			return nil
		}
	}
}

// UpdateManifest updates the key value entry for this service
// continuously, setting LastActive to the current Unix timestamp.
func (g *Gateway) UpdateManifest() {
	for {
		key := fmt.Sprintf("gateway/%s", g.ID)
		ts := time.Now().Unix()
		apps := []*consul.GatewayApp{
			&consul.GatewayApp{
				User:  "adamkdean",
				Name:  "hello-world",
				Image: "registry.dadi.engineer/adamkdean/hello-world",
			},
		}
		manifest := &consul.GatewayManifest{
			ID:         g.ID,
			Service:    "gateway",
			LastActive: ts,
			Apps:       apps,
		}
		fmt.Printf("Updating manifest, setting LastActive to %v\n", ts)
		if err := g.Consul.WriteStructToKey(key, manifest); err != nil {
			fmt.Printf("Error updating manifest: %v\n", err)
		}
		time.Sleep(time.Duration(g.UpdatePeriod) * time.Second)
	}
}

// Must handles errors and may include error reporting such
// as posting errors to a message queue before recovering.
func (g *Gateway) Must(err error) {
	if err != nil {
		// Log error? Recover?
		panic(err.Error())
	}
}

// New returns a new Gateway instance with the ID preset to
// an RFC4122 unique ID (See https://toolg.ietf.org/html/rfc4122).
func New() *Gateway {
	// Generate a UUID using V1 which incorporates both
	// timestamp and MAC address, and convert to string
	uuid := uuid.NewV1().String()

	return &Gateway{
		ID:           uuid,
		UpdatePeriod: 5,
	}
}
