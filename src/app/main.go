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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"lib"
	"log"
	"os"
	"sync"
)

// Application build version
var APP_VERSION string
var APP_NAME string

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool(
	"v", false, "Print the version number.")
var logfile *string = flag.String(
	"log", "file.log", "Log file to parse.")
var statemachine *string = flag.String(
	"state", "state.machine", "State machine file.")
var stateDescriptor *string = flag.String(
	"descriptor", "state.descriptor", "State description file")

// Initialize package name and version
func init() {
	if APP_VERSION == lib.EMPTY_STRING {
		APP_VERSION = "0.1"
	}
	if APP_NAME == lib.EMPTY_STRING {
		APP_NAME = "State Machine Checker"
	}
}

func main() {
	var errors []error
	var returnCode int
	returnCode = 0

	defer func() {
		if len(errors) != 0 {
			log.Printf("Error occurred: %v\n", errors)
		}
		os.Exit(returnCode)
	}()

	filesToRead := []string{"State Machine File (-state)",
		"State Description File (-descriptor)",
		"Log File (-log)"}
	fmt.Printf("%s, Version: %s\n", APP_NAME, APP_VERSION)
	fmt.Printf("The application utilizes %v\n", filesToRead)

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return
	}

	var err error

	// Read all the files, convert them to custom handlers
	var descriptorFile *os.File
	var descriptorFileReader *bufio.Scanner
	var statesFile *os.File
	var statesFileReader *bufio.Scanner
	var logFile *os.File
	var logFileReader *bufio.Scanner

	var wg sync.WaitGroup

	descriptorFile, err = os.Open(*stateDescriptor)
	if err != nil {
		returnCode = 1
		errors = append(errors, err)
		return
	}
	defer descriptorFile.Close()
	descriptorFileReader = bufio.NewScanner(descriptorFile)
	statesFile, err = os.Open(*statemachine)
	if err != nil {
		returnCode = 1
		errors = append(errors, err)
		return
	}
	defer statesFile.Close()
	statesFileReader = bufio.NewScanner(statesFile)
	logFile, err = os.Open(*logfile)
	if err != nil {
		returnCode = 1
		errors = append(errors, err)
		return
	}
	defer logFile.Close()
	logFileReader = bufio.NewScanner(logFile)

	// Read files simultaneously
	var stateDescription lib.StateDescription
	var stateMachine lib.StateMachine

	wg.Add(1)
	go func() {
		defer wg.Done()
		stateDescription, err = lib.BuildStateDescription(descriptorFileReader)
		if err != nil {
			errors = append(errors, err)
			returnCode = 2
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		stateMachine, err = lib.BuildStateMachine(statesFileReader)
		if err != nil {
			errors = append(errors, err)
			returnCode = 2
			return
		}
	}()
	wg.Wait()
	if returnCode != 0 {
		return
	}

	// Read the log file, line by line and make transitions
	var previousState map[string]lib.StateDefinition
	var state lib.StateDefinition
	var isAState bool

	// The number of elements in this map is the maximum number of
	// threads moving the state machine
	previousState = make(map[string]lib.StateDefinition)

	log.Println("Now parsing the log file")

	for {
		trace, isEnded, err := lib.ReadNextLine(logFileReader)
		if err != nil {
			errors = append(errors, err)
			returnCode = 3
			return
		}
		if isEnded {
			log.Println("Read the file completely")
			break
		}
		state, isAState, err = stateDescription.IdentifyState(trace)
		if err != nil {
			errors = append(errors, err)
			returnCode = 3
			return
		}
		if !isAState {
			// Just a log trace, go to next line
			// log.Printf("%v is not a state\n", trace)
			continue
		}

		prev, ok := previousState[state.Value]
		if !ok {
			// For this value, a first state, there's no state transition yet
			previousState[state.Value] = state
			continue
		}

		isCompleted, err := (&stateMachine).MakeTransition(prev, state)
		if err != nil {
			// Raise an exception, here's where to look for
			log.Printf("%v transitioned state from %v to %v\n", trace, prev, state)
			log.Printf("%v\n", err)
			log.Printf(
				"%v\nPlease refer to this found line for debugging", trace)
			return
		}
		if isCompleted {
			delete(previousState, state.Value)
		} else {
			previousState[state.Value] = state
		}
	}
	log.Println("Successfully processed the log file")
}
