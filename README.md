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
  mastodon    Simple Mastodon commands
  version     Print the version number

Flags:
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
  -h, --help                     help for toolbox
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox [command] --help" for more information about a command.
```

### Usage mastodon command

```
$ toolbox mastodon -h
Simple Mastodon commands.

Usage:
  toolbox mastodon [flags]
  toolbox mastodon [command]

Aliases:
  mastodon, mstdn, mast, mst

Available Commands:
  post        Post message to Mastodon
  profile     Output my profile
  register    Register application

Flags:
  -h, --help   help for mastodon

Global Flags:
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox mastodon [command] --help" for more information about a command.
```

### Usage bluesky command

```
$ toolbox bluesky -h
Simple Bluesky commands.

Usage:
  toolbox bluesky [flags]
  toolbox bluesky [command]

Aliases:
  bluesky, bsky, bs

Available Commands:
  post        Post message to Bluesky
  profile     Output Bluesky profile

Flags:
  -h, --help   help for bluesky

Global Flags:
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox bluesky [command] --help" for more information about a command.
```

## Modules Requirement Graph

[![dependency.png](./dependency.png)](./dependency.png)

[toolbox]: https://github.com/goark/toolbox "goark/toolbox: A collection of miscellaneous commands"
