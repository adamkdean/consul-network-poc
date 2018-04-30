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
	"github.com/adamkdean/consul-network-poc/utils/pkg/fsm"
	"github.com/adamkdean/consul-network-poc/utils/pkg/service"
	"github.com/adamkdean/consul-network-poc/utils/pkg/state"
	"github.com/satori/go.uuid"
	"time"
)

// Host is the application host layer for
// the DADI Cloud decentralized network.
type Host struct {
	ID           string
	Consul       *consul.Instance
	State        *fsm.StateMachine
	UpdatePeriod int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (h *Host) Initialize(consulAddr string) {
	defer h.Recover()
	h.InitializeState()
	h.InitializeService(consulAddr)
	h.InitializeManifestUpdateCycle()
	h.SearchForGateway()
}

// InitializeState creates a new state machine instance and
// hooks up an event to update the service state on change.
func (h *Host) InitializeState() {
	// Create a new state machine
	h.State = fsm.New()
	h.State.Initialize(map[string][]string{
		state.Initializing:        []string{state.SearchingForGateway},
		state.SearchingForGateway: []string{state.ConnectingToGateway, state.Error},
		state.ConnectingToGateway: []string{state.GatewayConnected, state.Error},
		state.GatewayConnected:    []string{},
		state.Error:               []string{},
	}, state.Initializing)

	// Update service status on state change.
	ch := make(chan string)
	go func() {
		for {
			st := <-ch
			fmt.Printf("Updating service to state %s\n", st)
			h.Must(h.UpdateService(st))
		}
	}()
	h.State.OnTransition("*", ch)
}

// InitializeService initializes the Consul client, then registers
// a service with them, and creates the current service manifest.
func (h *Host) InitializeService(consulAddr string) {
	h.Consul = consul.New()
	h.Must(h.Consul.Initialize(consulAddr))
	h.Must(h.UpdateService(h.State.CurrentState))
	h.Must(h.UpdateManifest())
}

// InitializeManifestUpdateCycle handles the manifest update cycle.
func (h *Host) InitializeManifestUpdateCycle() {
	go func() {
		for {
			if err := h.UpdateManifest(); err != nil {
				fmt.Printf("Error updating manifest: %v\n", err)
			}
			time.Sleep(time.Duration(h.UpdatePeriod) * time.Second)
		}
	}()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (h *Host) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := h.Consul.RegisterService(h.ID, service.Host, state); err != nil {
			fmt.Printf("Error registering service: %v (delay %v)\n", err, delay)
			if attempt > maxRetries {
				return fmt.Errorf("Could not register service")
			}

			time.Sleep(time.Duration(delay) * time.Second)
			delay *= 2
		} else {
			delay = 1
			attempt = 0
			fmt.Printf("Successfully registered service with ID %s and state %s\n", h.ID, state)
			return nil
		}
	}
}

// UpdateManifest updates the key value entry for this
// service, setting LastActive to the current Unix timestamp.
func (h *Host) UpdateManifest() error {
	ts := time.Now().Unix()
	key := fmt.Sprintf("%s/%s", service.Host, h.ID)
	manifest := &consul.ServiceManifest{
		ID:         h.ID,
		Type:       service.Host,
		LastActive: ts,
	}

	if err := h.Consul.WriteStructToKey(key, manifest); err != nil {
		return fmt.Errorf("error updating manifest: %v", err)
	}

	fmt.Printf("Updated manifest with LastActive %v\n", ts)
	return nil
}

// SearchForGateway finds all Gateway services with the
// state of AwaitingHosts, then tries to connect.
func (h *Host) SearchForGateway() {
	h.Must(h.State.Transition(state.SearchingForGateway))
	gateways, err := h.Consul.GetService(service.Gateway, state.AwaitingHosts)
	if err != nil {
		// TODO: wait x seconds, and then search again.
		fmt.Printf("Error searching for gateway: %v", err)
		h.Must(h.State.Transition(state.Error))
		return
	}

	// Iterate through Gateway services, attempting to connect.
	h.Must(h.State.Transition(state.ConnectingToGateway))
	for i := range gateways {
		g := gateways[i]
		fmt.Printf("Gateway found: %v\n", g.ID)
	}
}

// func (h *Host) ConnectToGateway() error {
// }

// Recover is used to recover from panic attacks.
func (h *Host) Recover() {
	if err := recover(); err != nil {
		fmt.Printf("Recovered from panic: %v\n", err)
	}
}

// Must handles errors and may include error reporting such
// as posting errors to a message queue before recovering.
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
	// timestamp and MAC address, and convert to string.
	uuid := uuid.NewV1().String()

	return &Host{
		ID:           uuid,
		UpdatePeriod: 5,
	}
}
