# [toolbox] -- A collection of miscellaneous commands

[![ci status](https://github.com/goark/toolbox/workflows/ci/badge.svg)](https://github.com/goark/toolbox/actions)
[![build status](https://github.com/goark/toolbox/workflows/build/badge.svg)](https://github.com/goark/toolbox/actions)
[![CodeQL status](https://github.com/goark/toolbox/workflows/CodeQL/badge.svg)](https://github.com/goark/toolbox/actions)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/goark/toolbox/master/LICENSE)
[![GitHub release](http://img.shields.io/github/release/goark/toolbox.svg)](https://github.com/goark/toolbox/releases/latest)

This package is required Go 1.27 or later.

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
  calendar    Astronomical calendar commands
  feed        Handling information for Web feed
  help        Help about any command
  mastodon    Simple Mastodon commands
  version     Print the version number
  webpage     Handling information for Web pages

Flags:
      --apod-config string       Config file for APOD (default "~/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "~/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "~/.cache/toolbox")
      --config string            Config file (default "~/.config/toolbox/config.yaml")
      --debug                    for debug
      --force-migration          Force database migration
  -h, --help                     help for toolbox
      --log-dir string           Directory for log files (default "~/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "~/.config/toolbox/mastodon.json")
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox [command] --help" for more information about a command.
```

### Usage mastodon command

```
$ toolbox mastodon -h
$ ./toolbox mastodon -h
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
      --apod-config string       Config file for APOD (default "~/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "~/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "~/.cache/toolbox")
      --config string            Config file (default "~/.config/toolbox/config.yaml")
      --debug                    for debug
      --force-migration          Force database migration
      --log-dir string           Directory for log files (default "~/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "~/.config/toolbox/mastodon.json")
      --temp-dir string          Temporary directory (default /tmp)

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
      --apod-config string       Config file for APOD (default "~/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "~/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "~/.cache/toolbox")
      --config string            Config file (default "~/.config/toolbox/config.yaml")
      --debug                    for debug
      --force-migration          Force database migration
      --log-dir string           Directory for log files (default "~/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "~/.config/toolbox/mastodon.json")
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox bluesky [command] --help" for more information about a command.
```

### Usage apod command

```
$ toolbox apod -h
Commands for Astronomy Picture of the Day by NASA/APOD API.

Usage:
  toolbox apod [flags]
  toolbox apod [command]

Available Commands:
  lookup      Lookup APOD data by NASA API
  post        Post APOD data to TL
  register    Register NASA API key

Flags:
  -d, --date string   Date for APOD data (YYYY-MM-DD)
  -h, --help          help for apod
  -u, --utc           Time base on UTC

Global Flags:
      --apod-config string       Config file for APOD (default "~/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "~/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "~/.cache/toolbox")
      --config string            Config file (default "~/.config/toolbox/config.yaml")
      --debug                    for debug
      --force-migration          Force database migration
      --log-dir string           Directory for log files (default "~/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "~/.config/toolbox/mastodon.json")
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox apod [command] --help" for more information about a command.
```

### APOD API migration notes (EN/JA)

Detailed bilingual notes are available here:
- [docs/apod-api-migration-notes.md](./docs/apod-api-migration-notes.md)

#### English

- Since 2026-09, APOD data source is switched to:
  - `https://science.nasa.gov/wp-json/wp/v2/apod-basic`
- Response format changed from NASA legacy APOD API to WordPress JSON API style.
- Current implementation behavior:
  - API requests use `page`/`per_page`.
  - `--date` and date-range semantics are handled by client-side filtering after fetch.
  - APOD metadata now uses fields like `permalink`, `credit`, `alt`, `url`, `hdurl`.
- SQLite schema note:
  - `apod_data` now includes `permalink` column.
  - Migration is executed even when database file already exists, so existing DBs are updated.
- Temporary behavior for social posting:
  - Credit text is sanitized by HTML tag stripping before message output.
  - If credit sanitization fails, the credit line is omitted and posting continues.

### Usage webpage command

```
$ toolbox webpage -h
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
      --save         Save page data to cache
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
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox webpage [command] --help" for more information about a command.
```

### Usage feed command

```
$ toolbox feed -h
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
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox feed [command] --help" for more information about a command.
```

### Usage calendar command

```
$ toolbox calendar -h
Commands for Astronomical calendar by NAOJ https://eco.mtk.nao.ac.jp/koyomi/cande/calendar.html.

Usage:
  toolbox calendar [flags]
  toolbox calendar [command]

Aliases:
  calendar, cal, c

Available Commands:
  lookup      Lookup astronomical calendar
  post        Post astronomical calendar data to TL

Flags:
      --eclipse           output eclipse
      --end string        end of date (YYYY-MM-DD)
      --ephemeris-all     output all ephemeris
  -h, --help              help for calendar
      --holiday           output holiday
      --moon-phase        output moon-phase
      --planet            output planet
      --solar-term        output solar-term
      --start string      start of date (YYYY-MM-DD)
      --template string   template file for Output format

Global Flags:
      --apod-config string       Config file for APOD (default "/home/username/.config/toolbox/nasaapi.json")
      --bluesky-config string    Config file for Bluesky (default "/home/username/.config/toolbox/bluesky.json")
      --cache-dir string         Directory for cache files (default "/home/username/.cache/toolbox")
      --config string            Config file (default "/home/username/.config/toolbox/config.yaml")
      --debug                    for debug
      --log-dir string           Directory for log files (default "/home/username/.cache/toolbox")
      --log-level string         Log level [nop|error|warn|info|debug|trace] (default "nop")
      --mastodon-config string   Config file for Mastodon (default "/home/username/.config/toolbox/mastodon.json")
      --temp-dir string          Temporary directory (default /tmp)

Use "toolbox calendar [command] --help" for more information about a command.
```

## Modules Requirement Graph

[![dependency.png](./dependency.png)](./dependency.png)

[toolbox]: https://github.com/goark/toolbox "goark/toolbox: A collection of miscellaneous commands"
