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

type StateDefinition struct {
	State string
	// If there's a variable/intelligent state machine log analysis requested
	// each state transition will be based on the value
	Value string
}

func NewStateDefinition(state string) StateDefinition {
	return StateDefinition{
		state,
		EMPTY_STRING,
	}
}

func NewEmptyStateDefinition() StateDefinition {
	return StateDefinition{
		EMPTY_STRING,
		EMPTY_STRING,
	}
}
