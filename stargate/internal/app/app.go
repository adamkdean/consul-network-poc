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

// Stargate is the authoritative DNS layer
// for the DADI Cloud decentralized network.
type Stargate struct {
	ID             string
	Consul         *consul.Instance
	State          *fsm.StateMachine
	UpdatePeriod   int
	PrunePeriod    int
	StaleThreshold int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (s *Stargate) Initialize(addr string) {
	defer s.Recover()
	s.InitializeState()
	s.InitializeService(addr)
	s.InitializeManifestUpdateCycle()
	s.InitializeServicePruneCycle()
}

// InitializeState creates a new state machine instance and
// hooks up an event to update the service state on change.
func (s *Stargate) InitializeState() {
	// Create a new state machine.
	s.State = fsm.New()
	s.State.Initialize(map[string][]string{
		state.Initializing: []string{state.Ready},
		state.Ready:        []string{},
	}, state.Initializing)

	// Update service status on state change.
	ch := make(chan string)
	go func() {
		for {
			st := <-ch
			fmt.Printf("Updating service to state %s\n", st)
			s.Must(s.UpdateService(st))
		}
	}()
	s.State.OnTransition("*", ch)
}

// InitializeService initializes the Consul client, then registers
// a service with them, and creates the current service manifest.
func (s *Stargate) InitializeService(addr string) {
	s.Consul = consul.New()
	s.Must(s.Consul.Initialize(addr))
	s.Must(s.UpdateService(s.State.CurrentState))
	s.Must(s.UpdateManifest())
	s.Must(s.State.Transition(state.Ready))
}

// InitializeManifestUpdateCycle handles the manifest update cycle.
func (s *Stargate) InitializeManifestUpdateCycle() {
	go func() {
		for {
			if err := s.UpdateManifest(); err != nil {
				fmt.Printf("Error updating manifest: %v\n", err)
			}
			time.Sleep(time.Duration(s.UpdatePeriod) * time.Second)
		}
	}()
}

// InitializeServicePruneCycle handles the service prune cycle.
func (s *Stargate) InitializeServicePruneCycle() {
	go func() {
		for {
			// Attempt to prune stale Hosts first.
			if err := s.PruneStaleServices(service.Host); err != nil {
				fmt.Printf("Error pruning stale Hosts: %v\n", err)
			}

			// Next, attempt to prune stale Gateaways.
			if err := s.PruneStaleServices(service.Gateway); err != nil {
				fmt.Printf("Error pruning stale Gateways: %v\n", err)
			}

			// Now we wait for the next cycle.
			time.Sleep(time.Duration(s.PrunePeriod) * time.Second)
		}
	}()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (s *Stargate) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := s.Consul.RegisterService(s.ID, service.Stargate, state); err != nil {
			fmt.Printf("Error registering service: %v (delay %v)\n", err, delay)
			if attempt > maxRetries {
				return fmt.Errorf("Could not register service")
			}

			time.Sleep(time.Duration(delay) * time.Second)
			delay *= 2
		} else {
			delay = 1
			attempt = 0
			fmt.Printf("Successfully registered service with ID %s and state %s\n", s.ID, state)
			return nil
		}
	}
}

// UpdateManifest updates the key value entry for this
// service, setting LastActive to the current Unix timestamp.
func (s *Stargate) UpdateManifest() error {
	ts := time.Now().Unix()
	key := fmt.Sprintf("%s/%s", service.Stargate, s.ID)
	manifest := &consul.ServiceManifest{
		ID:         s.ID,
		Type:       service.Stargate,
		LastActive: ts,
	}

	if err := s.Consul.WriteStructToKey(key, manifest); err != nil {
		return fmt.Errorf("error updating manifest: %v", err)
	}

	fmt.Printf("Updated manifest with LastActive %v\n", ts)
	return nil
}

// PruneStaleServices finds services that are inactive, and prunes them.
func (s *Stargate) PruneStaleServices(sv string) error {
	// Get a list of ServiceManifest for this service type.
	services, err := s.Consul.GetServiceManifests(sv)
	if err != nil {
		return err
	}

	// Debug variables.
	active := 0
	stale := 0
	removed := 0

	// What is the oldest the active timestamp can be?
	thr := time.Now().Add(time.Duration(s.StaleThreshold) * time.Second * -1)

	// Iterate through all services, and remove any that are stale.
	for i := range services {
		if thr.After(time.Unix(services[i].LastActive, 0)) {
			fmt.Printf("Service (%s) %s is stale, removing...\n", sv, services[i].ID)
			stale++

			// Deregister the service from Consul.
			if err := s.Consul.DeregisterService(services[i].ID); err != nil {
				return err
			}

			// Remove the manifest from the KV.
			if err := s.Consul.RemoveServiceManifest(sv, services[i].ID); err != nil {
				return err
			}

			fmt.Printf("Service (%s) %s removed\n", sv, services[i].ID)
			removed++
		} else {
			fmt.Printf("Service (%s) %s is active\n", sv, services[i].ID)
			active++
		}
	}

	fmt.Printf("Service (%s) total/active/stale/removed: %d/%d/%d/%d\n",
		sv, len(services), active, stale, removed)

	return nil
}

// Must handles errors and may include error reporting such
// as posting errors to a message queue before recovering.
func (s *Stargate) Must(err error) {
	if err != nil {
		// Log error?
		fmt.Printf("Panic! %v\n", err)
		panic(err.Error())
	}
}

// Recover is used to recover from panic attacks.
func (s *Stargate) Recover() {
	if err := recover(); err != nil {
		fmt.Printf("Recovered from panic: %v\n", err)
	}
}

// New returns a new Stargate instance with the ID preset to
// an RFC4122 unique ID (See https://tools.ietf.org/html/rfc4122).
func New() *Stargate {
	// Generate a UUID using V1 which incorporates both
	// timestamp and MAC address, and convert to string.
	uuid := uuid.NewV1().String()

	return &Stargate{
		ID:             uuid,
		UpdatePeriod:   5,
		PrunePeriod:    10,
		StaleThreshold: 20,
	}
}
