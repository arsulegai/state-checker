# State Machine Analysis Tool
This is a tool which analyzes possible state transitions in a given input log.
If the log traces differ from what the actual trace states, it will raise a
red flag.

The tool accepts three files as input `State Descriptor`, `State Machine` and
a `Log File`.

# Sample Configuration

The `State Descriptor` file looks like following

```
State1| Log Trace To Be Searched Which Corresponds to State 1
State2| Log Trace To Be Searched Which Corresponds to State 2
State3| Log Trace To Be Searched Which Corresponds to State 3
```

The `State Machine` file looks like following

```
State1|State2
State2|State3,State1
State3|State1
```

In above example, there's possible transition from State2 to State3 or State1.
But it should always go from State1 to State3 or State3 to State1.

The `Log` file can be any text file, example

```
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 1
[Log Trace Level, Time Information] Some Other Log Trace 1
[Log Trace Level, Time Information] Some Other Log Trace 2
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 2
[Log Trace Level, Time Information] Some Other Log Trace 3
[Log Trace Level, Time Information] Some Other Log Trace 4
[Log Trace Level, Time Information] Some Other Log Trace 5
[Log Trace Level, Time Information] Some Other Log Trace 6
[Log Trace Level, Time Information] Log Trace To Be Searched Which Corresponds to State 3
```

If above example, there's no error in parsing the file.

# Build and Run
The program is tested on Ubuntu 18.04 LTS (Bionic). Also the docker build
options provided generates a binary to execute on the bionic machine.

## Bare Metal Build
Refer to (Golang/Go Repository)[https://github.com/golang/go] to know how to
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
directory

```
docker-compose -f docker/compose/build.yaml up
```
It will mount the local directory to the container and create a `bin` directory
where you'll find the executable binary.

**Note:**
If the build is successful, docker-compose up should exit with
status code 0.

## Run
To run the application use following command format, assumes that generated
binary is in the `$PATH`

```
state-machine -descriptor state-descriptor.file -state state-machine.file -log log.file
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


&copy; Copyright 2019, Intel Corporation
