package lib

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"log"
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
		// log.Printf("Log being checked is %v\n", value)
		if err != nil {
			return stateDescription, err
		}
		if result {
			// Intelligent value place holder will be at position 0 which is
			// value
			// log.Println("Found a tag")
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
		// log.Printf("Line for reading is %v and is matched %v with the result %v\n", lineForReading, matched, result)
		if matched {
			toReturnState, ok := stateDescription.Description[description]
			// log.Printf("Raw state is %v\n", toReturnState)
			if !ok {
				return NewEmptyStateDefinition(), false, errors.New(fmt.Sprintf("Expected value not found %v", description))
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
				log.Printf("Left %v Right %v\n", leftPart, rightPart)
				matchedStateValue :=
					strings.TrimSpace(
						strings.Split(
							strings.Split(line, leftPart)[1], rightPart)[1])
				log.Printf("Matched %v\n", matchedStateValue)
				toReturnState.Value = matchedStateValue
			} else {
				toReturnState.Value = EMPTY_STRING
			}
			return toReturnState, true, nil
		}
	}
	return NewEmptyStateDefinition(), false, nil
}
