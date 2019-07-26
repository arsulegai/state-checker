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
