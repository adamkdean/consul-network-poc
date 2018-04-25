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
	"time"
	"github.com/adamkdean/consul-network-poc/utils/pkg/consul"
	"github.com/adamkdean/consul-network-poc/utils/pkg/state"
	"github.com/satori/go.uuid"
)

// Stargate is the authoritative DNS layer
// for the DADI Cloud decentralized network
type Stargate struct {
	ID     string
	Consul *consul.Instance
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (s *Stargate) Initialize(addr string) {
	s.Consul = consul.New()
	s.Must(s.Consul.Initialize(addr))
	s.UpdateService(state.Initialized)
	go s.UpdateManifest()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (s *Stargate) UpdateService(state string) {
	s.Must(s.Consul.RegisterService(s.ID, "stargate", state))
}

// UpdateManifest updates the key value entry for this service
// continuously, setting LastActive to the current Unix timestamp.
func (s *Stargate) UpdateManifest() {
	for {
		key := fmt.Sprintf("stargate/%s", s.ID)
		ts := time.Now().Unix()
		manifest := &consul.BasicManifest{
			ID:         s.ID,
			Service:    "stargate",
			LastActive: ts,
		}
		s.Must(s.Consul.WriteStructToKey(key, manifest))
		time.Sleep(1 * time.Second)
	}
}

// Must handles errors and may include error reporting such
// as posting errors to a message queue before recovering.
func (s *Stargate) Must(err error) {
	if err != nil {
		// Log error? Recover?
		panic(err.Error())
	}
}

// New returns a new Stargate instance with the ID preset to
// an RFC4122 unique ID (See https://tools.ietf.org/html/rfc4122).
func New() *Stargate {
	// Generate a UUID using V1 which incorporates both
	// timestamp and MAC address, and convert to string
	uuid := fmt.Sprintf("%s", uuid.Must(uuid.NewV1()))

	return &Stargate{
		ID: uuid,
	}
}
