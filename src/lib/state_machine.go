package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

func BuildStateMachine(fileReader *bufio.Scanner) (map[string][]string, error) {
	var stateMachine map[string][]string
	stateMachine = make(map[string][]string)
	const numberOfParts int = 2

	for {
		line, isEnded, err := ReadNextLine(fileReader)
		if err != nil {
			return nil, err
		}
		if isEnded {
			break
		}
		parts := strings.SplitN(line, STATE_DELIMITER, numberOfParts)
		if len(parts) != numberOfParts {
			return nil, errors.New(
				fmt.Sprintf(
					"Line %s has %d parts, but expected %d",
					line,
					len(parts),
					numberOfParts),
			)
		}
		possibleStates := strings.Split(parts[1], STATE_DELIMITER)
		stateMachine[parts[0]] = possibleStates
	}
	return stateMachine, nil
}

func MakeTransition(
	curState string,
	nextState string,
	stateMachine map[string][]string,
) error {
	possibleStates, ok := stateMachine[curState]
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
