package lib

import (
	"bufio"
)

func ReadNextLine(fileReader *bufio.Scanner) (string, bool, error) {
	if fileReader.Scan() {
		return fileReader.Text(), false, fileReader.Err()
	}
	return EMPTY_STRING, true, nil
}
