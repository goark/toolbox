# [toolbox] -- A collection of miscellaneous commands

[![check vulns](https://github.com/goark/toolbox/workflows/vulns/badge.svg)](https://github.com/goark/toolbox/actions)
[![lint status](https://github.com/goark/toolbox/workflows/lint/badge.svg)](https://github.com/goark/toolbox/actions)
[![lint status](https://github.com/goark/toolbox/workflows/build/badge.svg)](https://github.com/goark/toolbox/actions)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/goark/toolbox/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/goark/toolbox.svg)](https://github.com/goark/toolbox/releases/latest)

This package is required Go 1.20 or later.

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
  apod        NASA APOD commands
  bluesky     Simple Bluesky commands
  help        Help about any command
  mastodon    Simple Mastodon commands
  version     Print the version number
  webpage     Handling information for Web pages

Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
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
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
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
  register    Register account in local PC

Flags:
  -h, --help   help for bluesky

Global Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox bluesky [command] --help" for more information about a command.
```

### Usage apod command

```
$ toolbox apod -h
Commands for Astronomy Picture of the Day by NASA API.

Usage:
  toolbox apod [flags]
  toolbox apod [command]

Available Commands:
  lookup      Lookup APOD data NASA API key
  post        Post APOD data to TL
  register    Register NASA API key

Flags:
  -d, --date string   Date for APOD data (YYYY-MM-DD)
  -h, --help          help for apod
  -u, --utc           Time base on UTC

Global Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox apod [command] --help" for more information about a command.
```

### Usage webpage command

```
$ toolbox apod -h
Handling information for Web pages.

Usage:
  toolbox webpage [flags]
  toolbox webpage [command]

Aliases:
  webpage, web, w, bookmark, book, bm

Available Commands:
  lookup      Lookup information for Web page
  post        Post Web page's information to TL

Flags:
  -h, --help         help for webpage
      --save         Save APOD data to cache
  -u, --url string   Web page URL

Global Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox webpage [command] --help" for more information about a command.
```

### Usage feed command

```
$ toolbox apod -h
Handling information for Web feed.

Usage:
  toolbox feed [flags]
  toolbox feed [command]

Aliases:
  feed, rss

Available Commands:
  lookup      Lookup information for Web page
  post        Post Web page's information to TL

Flags:
  -f, --feed-list-file string   path of Feed list file
      --flickr-id string        Flickr ID
  -h, --help                    help for feed
      --save                    Save webpage data to cache
  -u, --url string              Feed URL

Global Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")

Use "toolbox feed [command] --help" for more information about a command.
```

## Modules Requirement Graph

[![dependency.png](./dependency.png)](./dependency.png)

[toolbox]: https://github.com/goark/toolbox "goark/toolbox: A collection of miscellaneous commands"
