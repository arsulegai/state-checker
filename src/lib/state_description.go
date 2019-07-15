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
		line, isMoreToRead, err := ReadNextLine(fileReader)
		if err != nil {
			return nil, err
		}
		if !isMoreToRead {
			break
		}
		parts := strings.SplitN(line, STATE_DELIMITER, numberOfParts-1)
		if len(parts) != numberOfParts {
			return nil, errors.New(
				fmt.Sprintf(
					"Line has %d parts, but expected %d",
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
	state, ok := stateDescription[line]
	if !ok {
		return EMPTY_STRING, false, nil
	}
	return state, true, nil
}
