package fielder

import (
	"container/ring"
	"errors"
	"sync"
)

// will be an implementation of a state machine which does not use switches (switches at scale are hard to read)
// will be easier to read and develop on than a giant nested switch
// will be easily extendable by other objects/interfaces

// the machine itself will be represented by a ring data type (https://pkg.go.dev/container/ring)
// a ring is ideal for our purposes, because we want states to switch between each other, and a ring is a bit like a linked list that never ends
// the ring will contain items, each item is a state in the overall machine. the start state/node should be the first item in the ring. or otherwise marked somehow
// as "start" (bool true?)
// we will start at the start node, and then as we move around the state machine, we will need to keep track of our current position on the ring so that we know
// what our current state is. the transitions between ring items are defined by functions within the state item

type RingAddress *ring.Ring

// the ring id represents the index on the ring where the item is

type StateId string

var SameStateNoUpdate = errors.New("state has no transition, is terminal, no update")

// the state id represents how states refer to each other

type StateMachine struct {
	*ring.Ring
	mu                 *sync.RWMutex // protect our ring
	Start              StateId
	IdRingAddressCache map[StateId]RingAddress // the state machine will reference the states using their ring address, but the states refer to each other with state ids
	// so we map them to each other
	ValueCache map[StateId]StateValue // the StateValue is the value of each state. if a state machine has nodes 0, 1, 2, 3,  then each number is the value
	// of a state in the machine. if the state machine has nodes "first", "next", "then", "last", then each string is the value of a state in the machine

	// we parse the initial ring and create the value cache at instantiation. Because i want to protect our state machines and keep them simple, we will not allow
	// writing to the state machine once its created. if you need to change it, just create a new one with the states you want
}

func (sm *StateMachine) lookupValueCacheId(in StateValue, equals func(i, j StateValue) bool) StateId {
	for k, v := range sm.ValueCache {
		if equals(v, in) {
			return k
		}
	}
	return ""
}

// "in" is the current state value, "testData" is the data that will be tested by the conditional questions to determine next state
// "equals" is a function that allows us to compare values without knowing the exact type ahead of time
func (sm *StateMachine) ProcessInMachine(in StateValue, testData any, equals func(i, j StateValue) bool) (StateValue, error) {
	stateId := sm.lookupValueCacheId(in, equals)
	// evaluate state with id stateId
	currentAddr, ok := sm.IdRingAddressCache[stateId]
	if !ok {
		return nil, errors.New("id does not exist in machine")
	}
	if currentAddr == nil {
		return nil, errors.New("address does not exist in machine")
	}
	currentState, ok := currentAddr.Value.(State)
	if !ok {
		return nil, errors.New("error retrieving state")
	}
	nextId, err := currentState.EvaluateTransition(testData)
	if err != nil {
		return nil, err
	}
	if nextId == "" {
		return nil, errors.New("next id is empty")
	}
	if nextId == stateId {
		// we havent switched states, return the same
		return nil, SameStateNoUpdate
	}

	value, ok := sm.ValueCache[nextId]
	if !ok {
		return nil, errors.New("next id does not exist in machine")
	}
	if value == nil {
		return nil, errors.New("value is nil for next id")
	}
	return value, nil
}

func (sm *StateMachine) PopulateRing(in ...State) {
	// we already created the ring and make it length "len(in)" so we know we can iterate safely through the items in the ring
	// Initialize the ring with the states. create the id -> ring address cache and the id -> value cache
	for _, v := range in {
		sm.Ring.Value = v
		sm.IdRingAddressCache[v.Id] = sm.Ring
		sm.ValueCache[v.Id] = v.StateValue
		sm.Ring = sm.Ring.Next()
	}
}

func NewStateMachine(states ...State) *StateMachine {
	if len(states) == 0 {
		return &StateMachine{}
	}
	sm := &StateMachine{
		Ring:               ring.New(len(states)),
		mu:                 new(sync.RWMutex),
		Start:              states[0].Id, // start node should be the first state in the list
		IdRingAddressCache: make(map[StateId]RingAddress),
		ValueCache:         make(map[StateId]StateValue),
	}
	sm.PopulateRing(states...)
	return sm
}

type StateValue FieldValue

type State struct {
	Id         StateId      // since rings dont have a beginning or end, we give each state an id so that we can find it later
	Matches    []Transition // the different transitions from this current state
	StateValue              // what is the value at this state?
	Start      bool         // is this the start state for the machine?
	Terminal   bool         // is this an end state for the machine? (no more transitions are needed)
}

func (s *State) EvaluateTransition(dataToTest any) (StateId, error) {
	if s.Terminal {
		return s.Id, nil
	}
	for _, v := range s.Matches {
		if v.SimpleMatcher(dataToTest) {
			return v.NextState, nil
		}
	}
	return "", errors.New("no valid transitions available")
}

// note, each transition matcher must be mutually exclusive
// only one transition of a state should evaluate to "true" at a time, there should be no overlap between transitions
type Transition struct {
	NextState     StateId
	SimpleMatcher // matcher is a much simpler function type to support whether we are eligible to move to this next state
}

type SimpleMatcher func(inputToMatch any) bool

func BasicEquals(s1, s2 StateValue) bool {
	return s1 == s2
}

func NextBehavior[inputType any, behaviorType any](sm *StateMachine, currentStatus StateValue, dataToTest inputType, mapper map[StateValue]behaviorType) (StateValue, behaviorType, error) {
	nextValue, err := sm.ProcessInMachine(currentStatus, dataToTest, BasicEquals)
	if err != nil {
		return "", *new(behaviorType), err
	}
	behavior, ok := mapper[nextValue]
	if !ok {
		return "", *new(behaviorType), errors.New("no mapped behavior available for state value")
	}
	return nextValue, behavior, nil
}
