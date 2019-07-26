package lib

import ()

const (
	STATE_UNKNOWN = iota
	KNOWN_STATE
)

const (
	EMPTY_STRING    string = ""
	WORD_DELIMITER  string = " "
	STATE_DELIMITER string = "|"
	LIST_DELIMITER  string = ","
	TRACE_DELIMITER string = "]"
	START_TAG       string = "<Value>"
	END_TAG         string = "</Value>"
	TAG_STRING      string = START_TAG + ".*" + END_TAG
	END_STATE       string = "END_STATE"
)
