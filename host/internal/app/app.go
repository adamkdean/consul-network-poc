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

// Host is the application host layer for
// the DADI Cloud decentralized network
type Host struct {
	ID           string
	Consul       *consul.Instance
	UpdatePeriod int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (h *Host) Initialize(addr string) {
	h.Consul = consul.New()
	h.Must(h.Consul.Initialize(addr))
	h.Must(h.UpdateService(state.Initialized))
	go h.UpdateManifest()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (h *Host) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := h.Consul.RegisterService(h.ID, "host", state); err != nil {
			fmt.Printf("Error registering service: %v (delay %v)\n", err, delay)
			if attempt > maxRetries {
				return fmt.Errorf("Could not register service")
			}

			time.Sleep(time.Duration(delay) * time.Second)
			delay *= 2
		} else {
			fmt.Printf("Successfully registered service with ID %s and state %s\n", h.ID, state)
			return nil
		}
	}
}

// UpdateManifest updates the key value entry for this service
// continuously, setting LastActive to the current Unix timestamp.
func (h *Host) UpdateManifest() {
	for {
		key := fmt.Sprintf("host/%s", h.ID)
		ts := time.Now().Unix()
		manifest := &consul.BasicManifest{
			ID:         h.ID,
			Service:    "host",
			LastActive: ts,
		}
		fmt.Printf("Updating manifest, setting LastActive to %v\n", ts)
		if err := h.Consul.WriteStructToKey(key, manifest); err != nil {
			fmt.Printf("Error updating manifest: %v\n", err)
		}
		time.Sleep(time.Duration(h.UpdatePeriod) * time.Second)
	}
}

// Must handles errors and may include error reporting such
// as posting errors to a message queue before recoverinh.
func (h *Host) Must(err error) {
	if err != nil {
		// Log error? Recover?
		panic(err.Error())
	}
}

// New returns a new Host instance with the ID preset to
// an RFC4122 unique ID (See https://toolh.ietf.org/html/rfc4122).
func New() *Host {
	// Generate a UUID using V1 which incorporates both
	// timestamp and MAC address, and convert to string
	uuid := uuid.NewV1().String()

	return &Host{
		ID:           uuid,
		UpdatePeriod: 5,
	}
}
