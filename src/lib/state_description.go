package lib

import (
	"bufio"
	"errors"
	"fmt"
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
		stateDescription[parts[1]] = parts[0]
	}
	return stateDescription, nil
}

func IdentifyState(
	line string,
	stateDescription map[string]string,
) (string, bool, error) {
	// State Description is straightaway a key in this case
	parts := strings.Split(line, TRACE_DELIMITER)
	trace := parts[len(parts) - 1]
	state, ok := stateDescription[trace]
	if !ok {
		return EMPTY_STRING, false, nil
	}
	return state, true, nil
}
