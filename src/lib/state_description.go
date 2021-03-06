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
	"regexp"
	"strings"
)

type StateDescription struct {
	Description map[string]StateDefinition
	Values      map[string]string
}

func newStateDescription() StateDescription {
	return StateDescription{
		make(map[string]StateDefinition),
		make(map[string]string),
	}
}

func BuildStateDescription(
	fileReader *bufio.Scanner,
) (StateDescription, error) {
	stateDescription := newStateDescription()
	const numberOfParts int = 2

	for {
		line, isEnded, err := ReadNextLine(fileReader)
		if err != nil {
			return stateDescription, err
		}
		if isEnded {
			break
		}
		parts := strings.SplitN(line, STATE_DELIMITER, numberOfParts)
		if len(parts) != numberOfParts {
			return stateDescription, errors.New(
				fmt.Sprintf(
					"Line %s has %d parts, but expected %d",
					line,
					len(parts),
					numberOfParts),
			)
		}
		key := strings.TrimSpace(parts[1])
		value := strings.TrimSpace(parts[0])
		result, err := regexp.MatchString(TAG_STRING, value)
		if err != nil {
			return stateDescription, err
		}
		if result {
			// Intelligent value place holder will be at position 0 which is
			// value
			extracted := strings.Split(
				strings.Split(value, END_TAG)[0], START_TAG)[1]
			stateDescription.Values[extracted] = key
		} else {
			stateDescription.Description[key] = NewStateDefinition(value)
		}
	}

	return stateDescription, nil
}

func (stateDescription StateDescription) IdentifyState(
	line string,
) (StateDefinition, bool, error) {
	// For each of the state description, check if the given line matches it
	for description := range stateDescription.Description {
		var lineForReading string
		result, err := regexp.MatchString(TAG_STRING, description)
		if err != nil {
			return NewEmptyStateDefinition(), false, err
		}
		if result {
			extracted := strings.Split(strings.Split(
				description, END_TAG)[0], START_TAG)[1]
			value, ok := stateDescription.Values[extracted]
			if !ok {
				return NewEmptyStateDefinition(),
					false,
					errors.New(
						"Unexpected error while trying to fetch the known value")
			}
			toReplace := START_TAG + extracted + END_TAG
			lineForReading = strings.Replace(description, toReplace, value, 1)
		} else {
			lineForReading = description
		}
		matched, err := regexp.MatchString(lineForReading, line)
		if err != nil {
			return NewEmptyStateDefinition(), false, err
		}
		if matched {
			toReturnState, ok := stateDescription.Description[description]
			if !ok {
				return NewEmptyStateDefinition(),
					false,
					errors.New(
						fmt.Sprintf("Expected value not found %v", description))
			}
			if result {
				leftPart := strings.TrimSpace(
					strings.Split(description, START_TAG)[0])
				if leftPart == EMPTY_STRING {
					leftPart = WORD_DELIMITER
				}
				rightPart := strings.TrimSpace(
					strings.Split(description, END_TAG)[1])
				if rightPart == EMPTY_STRING {
					rightPart = WORD_DELIMITER
				}
				matchedStateValue :=
					strings.TrimSpace(
						strings.Split(
							strings.Split(line, leftPart)[1], rightPart)[1])
				toReturnState.Value = matchedStateValue
			} else {
				toReturnState.Value = EMPTY_STRING
			}
			return toReturnState, true, nil
		}
	}
	return NewEmptyStateDefinition(), false, nil
}
