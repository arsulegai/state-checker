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
	if APP_VERSION == "" {
		APP_VERSION = "0.1"
	}
	if APP_NAME == "" {
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

	filesToRead := []string{"State Machine", "State Description", "Log File"}
	log.Printf("%s, Version: %s\n", APP_NAME, APP_VERSION)
	log.Printf("The application utilizes %v", filesToRead)

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		return
	}

	var err error

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

	var stateDescription map[string]string
	var stateMachine map[string][]string

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

	var previousState string
	var state string
	var isAState bool

	for {
		trace, isMoreRead, err := lib.ReadNextLine(logFileReader)
		if err != nil {
			errors = append(errors, err)
			returnCode = 3
			return
		}
		if !isMoreRead {
			break
		}
		previousState = state
		state, isAState, err = lib.IdentifyState(trace, stateDescription)
		if err != nil {
			errors = append(errors, err)
			returnCode = 3
			return
		}
		if !isAState {
			// Just a log trace, go to next line
			continue
		}

		err = lib.MakeTransition(previousState, state, stateMachine)
		if err != nil {
			// Raise exception, here's where to look for
			log.Printf("%v\n", err)
			log.Printf(
				"%s\nPlease refer to this found line for debugging", trace)
			return
		}
	}
}