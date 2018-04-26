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
	"github.com/adamkdean/consul-network-poc/utils/pkg/state"
	"github.com/satori/go.uuid"
	"time"
)

// Stargate is the authoritative DNS layer
// for the DADI Cloud decentralized network.
type Stargate struct {
	ID           string
	Consul       *consul.Instance
	State        *fsm.StateMachine
	UpdatePeriod int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (s *Stargate) Initialize(addr string) {
	defer s.Recover()
	s.InitializeState()
	s.InitializeService(addr)
	s.InitializeManifestUpdateCycle()
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

// InitializeService ...
func (s *Stargate) InitializeService(addr string) {
	s.Consul = consul.New()
	s.Must(s.Consul.Initialize(addr))
	s.Must(s.UpdateService(s.State.CurrentState))
	s.Must(s.UpdateManifest())
	s.Must(s.State.Transition(state.Ready))
}

// InitializeManifestUpdateCycle ...
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

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (s *Stargate) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := s.Consul.RegisterService(s.ID, "stargate", state); err != nil {
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
	key := fmt.Sprintf("stargate/%s", s.ID)
	manifest := &consul.StargateManifest{
		ID:         s.ID,
		Service:    "stargate",
		LastActive: ts,
	}

	if err := s.Consul.WriteStructToKey(key, manifest); err != nil {
		return fmt.Errorf("error updating manifest: %v", err)
	}

	fmt.Printf("Updated manifest with LastActive %v\n", ts)
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
		ID:           uuid,
		UpdatePeriod: 5,
	}
}
