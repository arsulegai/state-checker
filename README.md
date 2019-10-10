# State Machine Analysis Tool
This is a tool which analyzes possible state transitions in a given input log.
If the log traces differ from what the actual trace states, it will raise a
red flag.

The tool accepts three files as input `State Descriptor`, `State Machine` and
a `Log File`. The document further assumes the name of the descriptor file
to be `state-descriptor.file`, name of the state machine file to be
`state-machine.file` and the name of the log file to be `log.file`.

The tool use cases can be plenty, with the advancement in the tool there's
no limit to what one can do. However these can be widely used for following,
(Example output of these are given later in the document)

* Use Case 1: Catch the abnormality in state transition of the application,
with the ability to match log traces using regular expression (Go's `regexp`
is accepted).
* Use Case 2: Catch the abnormality in state transition for particular value.
This is helpful where the application makes state transitions for multiple
values simultaneously. For example, a web server serving request to multiple
clients at the same time using different session ids. The tool can be used
to identify if all the requests are handled gracefully as expected.

# Sample Configuration

The `state-descriptor.file` file looks like following

```
State1| Log Trace To Be Searched Which Corresponds to State 1, has a pattern [a-z0-9]+
State2| Log Trace To Be Searched Which Corresponds to State 2, has no pattern
State3| Log Trace To Be Searched Which Corresponds to State 3, has a pattern [a-z][0-9][A-Z]+
```
The `|` acts as `OR` operation. The line has to be read as `State1` or
`Log Trace To Be Searched Which Corresponds to State 1, has a pattern [a-z0-9]+`
both of them are considered equivalent. Henceforth the log trace will be referred
to using the state short-name in the `state-machine.file`.

**Note:** The descriptor can be a `Golang` [regexp](https://godoc.org/regexp)
expression. If regex is used, a log trace matching regex will be considered
as a state. Example is as follows

```
State1| Log Trace Pattern To Be Searched For <Value>VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE</Value>
State2| Another Trace Pattern For <Value>VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE</Value>
State3| Yet Another Trace Pattern For <Value>VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE</Value>
<Value>VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE</Value>| has a pattern [XYZ]
```
In this file, `<Value>` and `</Value>` are special tags. The log patterns matching
`VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE`'s defined pattern is identified.
This transition happens for each of the value matching
`VALUE-TO-BE-FOUND-FOR-TRANSITION-OF-STATE`, red flag is raised if any of the
value does not transition as expected.

**Note:** Different values can be overlapping in the log file. That is,
there can be multiple values performing state transition at the same time.
Tool handles them all.

The `state-machine.file` file looks like following

```
State1|State2
State2|State3,State1
State3|State1
```

In the above example, there's possible transition from State2 to State3 or
State1. The precedence of transition is that State2 will first try to move to
State3 before moving to State1. But it should always go from State1 to State2
or State3 to State1.

A special `END_STATE` can be used to stop the state transition abruptly. It's
helpful in positive manner for the `Use Case 2` which is descrived initially.
Example usage is discussed later in the document.

The `log.file` file can be any text file, example

```
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 1, has a pattern ab1234
[Log Trace Level, Time Information] Some Other Log Trace 1
[Log Trace Level, Time Information] Some Other Log Trace 2
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 2, has no pattern
[Log Trace Level, Time Information] Some Other Log Trace 3
[Log Trace Level, Time Information] Some Other Log Trace 4
[Log Trace Level, Time Information] Some Other Log Trace 5
[Log Trace Level, Time Information] Some Other Log Trace 6
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 3, has a pattern a1A
```

In the above example, there's no error in parsing the file.

The tool works as follows
1. Reads the input log file line by line
2. If any of the log line matches a state as per the descriptor file, it's
considered as a start of state-machine analysis.
3. If a log trace matches a state from descriptor file and transitions to
unexpected state (not as per the state-machine file), a red flag is raised.

**Note:** The state machine has no start or end. Because log file can be
sharded and the state machine can repeat itself in a loop. The way it works
is identifying first state in the log.

# Example Output

## Use Case 1: (Simple state transition with pattern matching)
`state-machine.file` is as follows

```
ConsensusNewMessage|ConsensusBlockValid
ConsensusBlockValid|StartCommitting,IgnoreBlock
StartCommitting|ConsensusBlockCommit
IgnoreBlock|ConsensusNewMessage
ConsensusBlockCommit|ConsensusNewMessage
```

`state-descriptor.file` is as follows
```
ConsensusNewMessage| Received message: CONSENSUS_NOTIFY_BLOCK_NEW
ConsensusBlockValid| Received message: CONSENSUS_NOTIFY_BLOCK_VALID
ConsensusBlockCommit| Received message: CONSENSUS_NOTIFY_BLOCK_COMMIT
StartCommitting| Committing [a-z0-9]+
IgnoreBlock| Ignoring [a-z0-9]+
```

`log.file` is a sample debug log file from
[Hyperledger Sawtooth PoET](https://github.com/hyperledger/sawtooth-poet)
application.
Output from the tool when run with the sample log file which has error is as follows

```
State Machine Checker, Version: 0.1
The application utilizes [State Machine File (-state) State Description File (-descriptor) Log File (-log)]
2019/07/26 21:55:07 Now parsing the log file
2019/07/26 21:55:07 [03:25:00.506 [MainThread] engine DEBUG] Received message: CONSENSUS_NOTIFY_BLOCK_NEW transitioned state from {StartCommitting } to {ConsensusNewMessage }
2019/07/26 21:55:07 Cannot transition from {StartCommitting } to {ConsensusNewMessage }
Possible are [{ConsensusBlockCommit }]

2019/07/26 21:55:07 [03:25:00.506 [MainThread] engine DEBUG] Received message: CONSENSUS_NOTIFY_BLOCK_NEW
Please refer to this found line for debugging
```

## Use Case 2: (Value based state transition)
`state-machine.file` is as follows

```
ConsensusBlockValid|CommitBlock,IgnoreBlock
CommitBlock|END_STATE
IgnoreBlock|END_STATE
END_STATE|ConsensusBlockValid
```

`state-descriptor.file` is as follows

```
ConsensusBlockValid| Passed consensus check: <Value>BLOCK_HERE</Value>
CommitBlock| Committing <Value>BLOCK_HERE</Value>
IgnoreBlock| Ignoring <Value>BLOCK_HERE</Value>
<Value>BLOCK_HERE</Value>| [a-z0-9]+
```

`log.file` is again a sample debug log file from
[Hyperledger Sawtooth PoET](https://github.com/hyperledger/sawtooth-poet)
application.
Output from the tool when run with the sample log file which has error is as follows

```
State Machine Checker, Version: 0.1
The application utilizes [State Machine File (-state) State Description File (-descriptor) Log File (-log)]
2019/07/26 22:03:10 Now parsing the log file
2019/07/26 22:03:10 [03:31:10.852 [MainThread] engine INFO] Passed consensus check: bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6 transitioned state from {ConsensusBlockValid bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6} to {ConsensusBlockValid bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6}
2019/07/26 22:03:10 Cannot transition from {ConsensusBlockValid bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6} to {ConsensusBlockValid bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6}
Possible are [{CommitBlock } {IgnoreBlock }]

2019/07/26 22:03:10 [03:31:10.852 [MainThread] engine INFO] Passed consensus check: bb3e98cfd3de79fce125a62d76a6cce788b8ea1357b0400a79295cd4e5d10cbf2f429fdc4cf03c9859f4e231243f82796f46f44934cf8a3c8556d0bfa821b6f6
Please refer to this found line for debugging
```

# Build and Run
The program is tested on Ubuntu 18.04 LTS (Bionic). Also the docker build
options provided generates a binary to execute on the bionic machine.

## Bare Metal Build
Refer to [Golang/Go Repository](https://github.com/golang/go) to know how to
configure Go in your machine. To build the binary of the standalone application
run the following command from within `src/app` directory. Also set project
root directory in `$GOPATH`.

```
go build
```

## Docker Build
**Note:**
Tested on
1. Docker-compose version 1.22
2. Docker version 18.06-ce

To generate a binary, run the docker-compose file from within the root
directory of the repository (example, where you clone the git repository)

```
docker-compose -f docker/compose/build.yaml up
```
It will mount the local directory to the container and create a `bin` directory
where you'll find the executable binary.

**Note:**
- If the build is successful, docker-compose up should exit with
status code 0.
- For generating a binary for the Mac. Please use the file `docker/compose/mac-build.yaml`.

## Run
To run the application use following command format, assumes that generated
binary is in the `$PATH`. If you are using the docker compose file for generation
then add `<full-path-to-repository-root-folder>/bin` in `$PATH`.

```
state-machine-analyzer -descriptor state-descriptor.file -state state-machine.file -log log.file
```

## Help

```
state-machine-analyzer --help
```

```
State Machine Checker, Version: 0.1
The application utilizes [State Machine File (-state) State Description File (-descriptor) Log File (-log)]
Usage of state-machine-analyzer:
  -descriptor string
        State description file (default "state.descriptor")
  -log string
        Log file to parse. (default "file.log")
  -state string
        State machine file. (default "state.machine")
  -v    Print the version number.
```

## Machine in proxy network
If you're using the tool on a machine in a proxy network environment, the build
may fail to get required packages or may fail during the docker image creation.
Please create a file `config.json` with following contents and place it under
the `/home/$USER/.docker/` directory. Create the directory if not present
already. The file looks like the following

```
{
 "proxies":
 {
   "default":
   {
     "httpProxy": "http://proxy-address-here:<proxy-port-http>",
     "httpsProxy": "http://proxy-address-here:<proxy-port-https>",
     "noProxy": "127.0.0.1,localhost",
     "hkpProxy": "http://proxy-address-here:<proxy-port-hkp>"
   }
 }
}
```

# Developers

## Contributions
You're free to improvise the application, raise a pull request to the original
repository after your implementation. Each commit must include `Signed-off-by:`
in the commit message (run `git commit -s` to auto-sign). This sign off means
you agree the commit satisfies the [Developer Certificate of
Origin(DCO)](https://developercertificate.org/).

## Beautiful Go
For the benefit of new code gazers, run the `go fmt` before raising the pull
request to the [https://github.com/arsulegai/state-checker](GitHub). There's
a docker compose file for help as well. Run the command from root directory
of the repository.

```
docker-compose -f docker/compose/fmt.yaml up
```
**Note:** Command will exit with the code 0 upon success.

## License
This software is licensed under the [Apache License Version 2.0](LICENSE)
software license.

&copy; Copyright 2019, Intel Corporation
