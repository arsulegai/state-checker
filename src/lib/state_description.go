package lib

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func BuildStateDescription(
	fileReader *bufio.Scanner,
) (map[string]string, error) {
	var stateDescription map[string]string
	stateDescription = make(map[string]string)
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
		key := strings.TrimSpace(parts[1])
		value := strings.TrimSpace(parts[0])
		stateDescription[key] = value
	}
	return stateDescription, nil
}

func IdentifyState(
	line string,
	stateDescription map[string]string,
) (string, bool, error) {
	// For each of the state description, check if the given line matches it
	for description := range stateDescription {
		matched, err := regexp.MatchString(description, line)
		if matched {
			return stateDescription[description], true, nil
		}
		if err != nil {
			return EMPTY_STRING, false, err
		}
	}
	return EMPTY_STRING, false, nil
}
