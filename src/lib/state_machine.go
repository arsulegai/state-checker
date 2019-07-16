package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type StateMachine struct {
	Machine map[string][]string
}

func newStateMachine() StateMachine {
	return StateMachine{
		make(map[string][]string),
	}
}

func BuildStateMachine(fileReader *bufio.Scanner) (StateMachine, error) {
	stateMachine := newStateMachine()
	const numberOfParts int = 2

	for {
		line, isEnded, err := ReadNextLine(fileReader)
		if err != nil {
			return stateMachine, err
		}
		if isEnded {
			break
		}
		parts := strings.SplitN(line, STATE_DELIMITER, numberOfParts)
		if len(parts) != numberOfParts {
			return stateMachine, errors.New(
				fmt.Sprintf(
					"Line %s has %d parts, but expected %d",
					line,
					len(parts),
					numberOfParts),
			)
		}
		initialState := strings.TrimSpace(parts[0])
		possibleStates := strings.Split(parts[1], LIST_DELIMITER)
		finalStates := []string{}
		for _, possibleState := range possibleStates {
			finalStates = append(finalStates, possibleState)
		}
		stateMachine.Machine[initialState] = finalStates
	}
	return stateMachine, nil
}

func (stateMachine *StateMachine) MakeTransition(
	curState string,
	nextState string,
) error {
	possibleStates, ok := stateMachine.Machine[curState]
	if !ok {
		// Nothing to do with this state
		return errors.New(
			fmt.Sprintf("Undefined error, no path found for %s", curState))
	}
	for _, state := range possibleStates {
		if nextState == state {
			// Found possible transition
			return nil
		}
	}
	// Current State to Next State couldn't be transitioned
	return errors.New(
		fmt.Sprintf(
			"Cannot transition from %s to %s\nPossible are %v",
			curState,
			nextState,
			possibleStates,
		),
	)
}
