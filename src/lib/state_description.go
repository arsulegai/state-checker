package lib

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type StateDescription struct {
	Description map[string]string
	Values      []string
}

func newStateDescription() StateDescription {
	return StateDescription{
		make(map[string]string),
		[]string{},
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
		stateDescription.Description[key] = value
	}
	return stateDescription, nil
}

func (stateDescription StateDescription) IdentifyState(
	line string,
) (string, bool, error) {
	// For each of the state description, check if the given line matches it
	for description := range stateDescription.Description {
		matched, err := regexp.MatchString(description, line)
		if matched {
			return stateDescription.Description[description], true, nil
		}
		if err != nil {
			return EMPTY_STRING, false, err
		}
	}
	return EMPTY_STRING, false, nil
}
