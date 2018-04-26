//    ___                      _     ___  ___  ___
//   / __\___  _ __  ___ _   _| |   / _ \/___\/ __\
//  / /  / _ \| '_ \/ __| | | | |  / /_)//  // /
// / /__| (_) | | | \__ \ |_| | | / ___/ \_// /___
// \____/\___/|_| |_|___/\__,_|_| \/   \___/\____/
//
// Consul Network proof of concept
// (c) 2018 Adam K Dean

// Finite State Machine, based on
// https://github.com/adamkdean/fsm
package fsm

import (
	"fmt"
	"github.com/thoas/go-funk"
)

// StateMachine is the finite state machine struct
type StateMachine struct {
	CurrentState string
	States       []string
	StateMap     map[string][]string
	EventMap     map[string][]chan string
}

// Initialize takes a state map and an initial
// state and initializes the state machine
func (s *StateMachine) Initialize(sm map[string][]string, st string) {
	s.CurrentState = st
	s.States = funk.Keys(sm).([]string)
	s.StateMap = sm
	s.EventMap = map[string][]chan string{}
}

// Transition changes the state when permissable
func (s *StateMachine) Transition(to string) error {
	if err := s.assureStateExists(to); err != nil {
		return err
	}

	// Iterate through all valid transitions and ensure
	// the request transition state is allowed
	for _, st := range s.StateMap[s.CurrentState] {
		if st == to {
			s.CurrentState = to

			// Iterate through events for this new state
			for _, e := range s.EventMap[s.CurrentState] {
				e <- s.CurrentState
			}

			return nil
		}
	}

	return fmt.Errorf("Invalid transition: %v", to)
}

// OnTransition hooks up event channels to state transitions
func (s *StateMachine) OnTransition(st string, ch chan string) error {
	if err := s.assureStateExists(st); err != nil {
		return err
	}

	s.EventMap[st] = append(s.EventMap[st], ch)
	return nil
}

func (s *StateMachine) assureStateExists(st string) error {
	if !funk.Contains(s.States, st) {
		return fmt.Errorf("Invalid state: %v", st)
	}
	return nil
}

// New returns a new, empty StateMachine instance
func New() *StateMachine {
	return &StateMachine{}
}
