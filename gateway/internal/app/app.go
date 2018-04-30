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
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"time"
)

// Gateway is the application gateway layer
// for the DADI Cloud decentralized network.
type Gateway struct {
	ID           string
	Address      string
	Port         int
	Hosts        []string
	Apps         []*consul.App
	Consul       *consul.Instance
	Router       *mux.Router
	State        *fsm.StateMachine
	UpdatePeriod int
}

// Initialize the service, creating a new instance of
// Consul and updating the service & manifest loop.
func (g *Gateway) Initialize(consulAddr string, listenPort int) {
	defer g.Recover()
	g.InitializeState()
	g.InitializeService(consulAddr, listenPort)
	g.InitializeManifestUpdateCycle()
	g.InitializeWebServer()
}

// InitializeState creates a new state machine instance and
// hooks up an event to update the service state on change.
func (g *Gateway) InitializeState() {
	// Create a new state machine.
	g.State = fsm.New()
	g.State.Initialize(map[string][]string{
		state.Initializing:  []string{state.AwaitingHosts},
		state.AwaitingHosts: []string{state.Error},
		state.Error:         []string{},
	}, state.Initializing)

	// Update service status on state change.
	ch := make(chan string)
	go func() {
		for {
			st := <-ch
			fmt.Printf("Updating service to state %s\n", st)
			g.Must(g.UpdateService(st))
		}
	}()
	g.State.OnTransition("*", ch)
}

// InitializeService initializes the Consul client, then registers
// a service with them, and creates the current service manifest.
func (g *Gateway) InitializeService(consulAddr string, listenPort int) {
	g.Address, _ = os.Hostname()
	g.Port = listenPort
	g.Hosts = []string{}
	g.Apps = []*consul.App{
		&consul.App{
			User:  "adamkdean",
			Name:  "hello-world",
			Image: "registry.dadi.engineer/adamkdean/hello-world",
		},
	}
	g.Consul = consul.New()
	g.Must(g.Consul.Initialize(consulAddr))
	g.Must(g.UpdateService(g.State.CurrentState))
	g.Must(g.UpdateManifest())
}

// InitializeManifestUpdateCycle handles the manifest update cycle.
func (g *Gateway) InitializeManifestUpdateCycle() {
	go func() {
		for {
			if err := g.UpdateManifest(); err != nil {
				fmt.Printf("Error updating manifest: %v\n", err)
			}
			time.Sleep(time.Duration(g.UpdatePeriod) * time.Second)
		}
	}()
}

// UpdateService updates the current service within Consul
// with the state that is passed as the service "tag".
func (g *Gateway) UpdateService(state string) error {
	delay := 1
	attempt := 0
	maxRetries := 5

	for {
		attempt++
		if err := g.Consul.RegisterService(g.ID, service.Gateway, state); err != nil {
			fmt.Printf("Error registering service: %v (delay %v)\n", err, delay)
			if attempt > maxRetries {
				return fmt.Errorf("Could not register service")
			}

			time.Sleep(time.Duration(delay) * time.Second)
			delay *= 2
		} else {
			delay = 1
			attempt = 0
			fmt.Printf("Successfully registered service with ID %s and state %s\n", g.ID, state)
			return nil
		}
	}
}

// UpdateManifest updates the key value entry for this
// service, setting LastActive to the current Unix timestamp.
func (g *Gateway) UpdateManifest() error {
	ts := time.Now().Unix()
	key := fmt.Sprintf("%s/%s", service.Gateway, g.ID)
	manifest := &consul.ServiceManifest{
		ID:         g.ID,
		Type:       service.Gateway,
		Address:    g.Address,
		Port:       g.Port,
		Apps:       g.Apps,
		Hosts:      g.Hosts,
		LastActive: ts,
	}

	if err := g.Consul.WriteStructToKey(key, manifest); err != nil {
		return fmt.Errorf("error updating manifest: %v", err)
	}

	fmt.Printf("Updated manifest %v\n", manifest)
	return nil
}

// InitializeWebServer starts a simple HTTP server and
// listens for Hosts registering with the Gateway.
func (g *Gateway) InitializeWebServer() {
	addr := fmt.Sprintf("%s:%v", g.Address, g.Port)
	g.Router = mux.NewRouter()
	g.Router.HandleFunc("/register/{id}", g.OnRegister).Methods("POST")
	g.Must(g.State.Transition(state.AwaitingHosts))
	if err := http.ListenAndServe(addr, g.Router); err != nil {
		g.Must(g.State.Transition(state.Error))
	}
}

// OnRegister handles the POST /register/{id} route,
// adding new host's to the g.Host string array.
func (g *Gateway) OnRegister(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Make sure the Host ID isn't an empty string.
	if id == "" {
		fmt.Printf("[400 Bad Request] Host ID is null: %s\n", r.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure we haven't already registered this Host.
	for h := range g.Hosts {
		if g.Hosts[h] == id {
			fmt.Println("[302 Found] Host already registered with Gateway")
			w.WriteHeader(http.StatusFound)
			return
		}
	}
	fmt.Printf("[202 Accepted] Adding Host with ID %s\n", id)
	g.Hosts = append(g.Hosts, id)
	w.WriteHeader(http.StatusAccepted)
}

// Recover is used to recover from panic attacks.
func (g *Gateway) Recover() {
	if err := recover(); err != nil {
		fmt.Printf("Recovered from panic: %v\n", err)
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
	// timestamp and MAC address, and convert to string.
	uuid := uuid.NewV1().String()

	return &Gateway{
		ID:           uuid,
		UpdatePeriod: 5,
	}
}
