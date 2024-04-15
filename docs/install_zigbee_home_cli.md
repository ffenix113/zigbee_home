---
title: Install zigbee_home CLI
---

CLI is the main entry point for generating, building and flashing the firmware from one single place.

Currently installing from source is available. In later stages compiled binaries will be provided on Github.

[Go installation](https://go.dev/doc/install) is required to build CLI from source. Go version of `1.21.4` or later is required.
One might be already provided in your Linux distribution or `brew` on Mac.

## Installing using `go install`
This will be best for users that don't want to mess with source code or have no need to modify internal functionality.

To install the CLI just run
```bash
go install github.com/ffenix113/zigbee_home/cmd/zigbee@develop
```
This will install the CLI and it can be later used by running `zigbee` from command line.

!!! note
    The name of the executable is quite generic, if this is a problem - please see next runnning method.

## Running from source without installation
This solution is for more advanced users, that have some knowledge in git and Go.

To run an executable first you would need to pull the source code of the project to some directory:
```bash
$ git clone git@github.com:ffenix113/zigbee_home.git
```
Inside the cloned directory run
```bash
$ go run ./zigbee_home/cmd/zigbee/... [args...]
```

This will not add any executables in your PATH. Instead you would need to execute `go run` command mentioned above each time to run the CLI.

### Updating the source code
When using this method the code will not be updated automatically in any way. To do this you would need to navigate to cloned repository and run
```bash
$ git pull
```