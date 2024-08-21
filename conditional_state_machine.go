package fielder

import (
	"container/ring"
	"errors"
	"sync"
)

type ConditionalStateMachine struct {
	*StateMachine
}

func NewConditionalStateMachine(states ...ConditionalState) *ConditionalStateMachine {
	if len(states) == 0 {
		return &ConditionalStateMachine{StateMachine: &StateMachine{}}
	}
	sm := &ConditionalStateMachine{StateMachine: &StateMachine{
		Ring:               ring.New(len(states)),
		mu:                 new(sync.RWMutex),
		Start:              states[0].Id, // start node should be the first state in the list
		IdRingAddressCache: make(map[StateId]RingAddress),
		ValueCache:         make(map[StateId]StateValue),
	}}
	sm.PopulateConditionalRing(states...)
	return sm
}

func (sm *ConditionalStateMachine) PopulateConditionalRing(in ...ConditionalState) {
	// we already created the ring and make it length "len(in)" so we know we can iterate safely through the items in the ring
	// Initialize the ring with the states. create the id -> ring address cache and the id -> value cache
	for _, v := range in {
		sm.Ring.Value = v
		sm.IdRingAddressCache[v.Id] = sm.Ring
		sm.ValueCache[v.Id] = v.StateValue
		sm.Ring = sm.Ring.Next()
	}
}

// "in" is the current state value, "testData" is the data that will be tested by the conditional questions to determine next state
// "equals" is a function that allows us to compare values without knowing the exact type ahead of time
func (sm *ConditionalStateMachine) ProcessInMachine(in StateValue, testData any, equals func(i, j StateValue) bool) (StateValue, error) {
	stateId := sm.lookupValueCacheId(in, equals)
	// evaluate state with id stateId
	currentAddr, ok := sm.IdRingAddressCache[stateId]
	if !ok {
		return nil, errors.New("id does not exist in machine")
	}
	if currentAddr == nil {
		return nil, errors.New("address does not exist in machine")
	}
	currentState, ok := currentAddr.Value.(ConditionalState)
	if !ok {
		return nil, errors.New("error retrieving state")
	}
	//
	nextId, err := currentState.EvaluateTransition(testData)
	if err != nil {
		return nil, err
	}
	if nextId == "" {
		return nil, errors.New("next id is empty")
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

type ConditionalState struct {
	Id         StateId                 // since rings dont have a beginning or end, we give each state an id so that we can find it later
	Outcomes   []ConditionalTransition // the different transitions from this current state
	StateValue                         // what is the value at this state?
	Start      bool                    // is this the start state for the machine?
}

func (s *ConditionalState) EvaluateTransition(dataToTest any) (StateId, error) {
	for _, v := range s.Outcomes {
		if v.Conditional.Meets(dataToTest) {
			return v.NextState, nil
		}
	}
	return "", errors.New("no valid transitions available")
}

// note, each transition conditional must be mutually exclusive
// only one transition of a state should evaluate to "true" at a time, there should be no overlap between transitions
type ConditionalTransition struct {
	NextState   StateId
	Conditional // we use the conditional interface to support whether we are eligible to move to this next state
}

func example() {
	_ = NewConditionalStateMachine([]ConditionalState{{
		Id: "id1",
		Outcomes: []ConditionalTransition{{
			NextState: "id1",
			Conditional: Conditions([]Prerequisite{
				{
					IsCandidate: func(f any) bool {
						return false
					},
					Gauntlet: []Question{
						func() Enforceable {
							return func(f any) bool {
								return false
							}
						},
					},
				},
				{
					IsCandidate: func(f any) bool {
						return false
					},
					Gauntlet: []Question{
						func() Enforceable {
							return func(f any) bool {
								return false
							}
						},
					},
				},
			}...),
		}},
		StateValue: "",
		Start:      true,
	},
		{
			Id:         "id2",
			Outcomes:   nil,
			StateValue: nil,
			Start:      false,
		},
		{
			Id:         "id3",
			Outcomes:   nil,
			StateValue: nil,
			Start:      false,
		},
	}...)
}
