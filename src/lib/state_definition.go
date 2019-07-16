package lib

import ()

type StateDefinition struct {
	State string
	// If there's a variable/intelligent state machine log analysis requested
	// each state transition will be based on the value
	Value string
}

func NewStateDefinition(state string) StateDefinition {
	return StateDefinition{
		state,
		EMPTY_STRING,
	}
}

func NewEmptyStateDefinition() StateDefinition {
	return StateDefinition{
		EMPTY_STRING,
		EMPTY_STRING,
	}
}
