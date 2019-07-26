/**
 * Copyright 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package lib

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type StateMachine struct {
	Machine map[StateDefinition][]StateDefinition
}

func newStateMachine() StateMachine {
	return StateMachine{
		make(map[StateDefinition][]StateDefinition),
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
		initialState := NewStateDefinition(strings.TrimSpace(parts[0]))
		possibleStates := strings.Split(parts[1], LIST_DELIMITER)
		finalStates := []StateDefinition{}
		for _, possibleState := range possibleStates {
			finalStates =
				append(finalStates,
					NewStateDefinition(strings.TrimSpace(possibleState)))
		}
		stateMachine.Machine[initialState] = finalStates
	}
	return stateMachine, nil
}

func (stateMachine *StateMachine) MakeTransition(
	curState StateDefinition,
	nextState StateDefinition,
) (bool, error) {
	tempState := curState
	tempState.Value = EMPTY_STRING
	possibleStates, ok := stateMachine.Machine[tempState]
	if !ok {
		// Nothing to do with this state
		return false, errors.New(
			fmt.Sprintf("Undefined error, no path found for %s", curState))
	}
	for _, state := range possibleStates {
		if state.State == END_STATE {
			// Possible transition and it's end state for this Value
			return true, nil
		}
		if nextState.State == state.State {
			// Found possible transition
			return false, nil
		}
	}
	// Current State to Next State couldn't be transitioned
	return false, errors.New(
		fmt.Sprintf(
			"Cannot transition from %s to %s\nPossible are %v\n",
			curState,
			nextState,
			possibleStates,
		),
	)
}
