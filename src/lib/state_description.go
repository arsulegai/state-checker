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
			extracted := strings.TrimLeft(
				strings.TrimRight(value, END_TAG), START_TAG)
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
			extracted := strings.TrimLeft(strings.TrimRight(
				description, END_TAG), START_TAG)
			value, ok := stateDescription.Values[extracted]
			if !ok {
				return NewEmptyStateDefinition(),
					false,
					errors.New(
						"Unexpected error while trying to fetch the known value")
			}
			lineForReading = strings.Replace(description, TAG_STRING, value, 1)
		} else {
			lineForReading = description
		}
		matched, err := regexp.MatchString(lineForReading, line)
		if matched {
			toReturnState := stateDescription.Description[description]
			if result {
				leftPart := strings.Trim(
					strings.TrimRight(description, START_TAG), START_TAG)
				rightPart := strings.Trim(
					strings.TrimLeft(description, END_TAG), END_TAG)
				matchedStateValue :=
					strings.TrimSpace(
						strings.TrimLeft(
							strings.TrimRight(line, rightPart), leftPart))
				toReturnState.Value = matchedStateValue
			} else {
				toReturnState.Value = EMPTY_STRING
			}
			return toReturnState, true, nil
		}
		if err != nil {
			return NewEmptyStateDefinition(), false, err
		}
	}
	return NewEmptyStateDefinition(), false, nil
}
