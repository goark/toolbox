# [toolbox] -- A collection of miscellaneous commands

[![check vulns](https://github.com/goark/toolbox/workflows/vulns/badge.svg)](https://github.com/goark/toolbox/actions)
[![lint status](https://github.com/goark/toolbox/workflows/lint/badge.svg)](https://github.com/goark/toolbox/actions)
[![lint status](https://github.com/goark/toolbox/workflows/build/badge.svg)](https://github.com/goark/toolbox/actions)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/goark/toolbox/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/goark/toolbox.svg)](https://github.com/goark/toolbox/releases/latest)

This package is required Go 1.16 or later.

## Build and Install

```
$ go install github.com/goark/toolbox@latest
```

## Binaries

See [latest release](https://github.com/goark/toolbox/releases/latest).

## Usage

```
$ toolbox -h
A collection of miscellaneous commands.

Usage:
  toolbox [flags]
  toolbox [command]

Available Commands:
  bluesky     Simple Bluesky commands
  help        Help about any command
  version     Print the version number

Flags:
      --config string      Config file (default "/home/spiegel/.config/toolbox/config.yaml")
      --debug              for debug
  -h, --help               help for toolbox
      --log-dir string     Directory for log files (default "/home/spiegel/.cache/toolbox")
      --log-level string   Log level [nop|error|warn|info|debug|trace] (default "nop")

Use "toolbox [command] --help" for more information about a command.
```

## Modules Requirement Graph

[![dependency.png](./dependency.png)](./dependency.png)

[toolbox]: https://github.com/goark/toolbox "goark/toolbox: A collection of miscellaneous commands"
