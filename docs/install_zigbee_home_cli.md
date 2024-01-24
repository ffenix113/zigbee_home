---
title: Install zigbee_home CLI
---

CLI is the main entry point for generating, building and flashing the firmware from one single place.

Currently installing from source is available. In later stages compiled binaries will be provided on Github.

[Go installation](https://go.dev/doc/install) is required to build CLI from source. Go version of `1.21.4` or later is required.
One might be already provided in your Linux distribution or `brew` on Mac.

## Installing from source
To install the CLI just run
```bash
go install github.com/ffenix113/zigbee_home/cli/cmd/zigbee@develop
```
This will install the CLI and it can be later used by running `zigbee` from command line.

!!! note
    The name of the executable is quite generic, if this is a problem - please see next runnning method.

## Running from source without installation
To run an executable there are two steps:
1. Clone the repo: `git clone git@github.com:ffenix113/zigbee_home.git`
2. Execute `go run ./zigbee_home/cli/cmd/zigbee/... [args...]`

This will not add any executables in your PATH. Instead you would need to execute `go run` command mentioned above each time to run the CLI.