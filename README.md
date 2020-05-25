# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Run this as a cron job to periodically check custom data sources and automatically create Trello cards based on custom filters.

An example use case (which is already implemented) could be to create a Trello card for each GitHub issue that's been assigned to you.

## Currently Available Sources
 * Github Issues - https://github.com/issues/assigned
 * TodoDock Tasks - https://tododock.com

Feel free to add new sources or improve the implementations of the existing ones. Contributions are always welcome!

## Configuration
Copy and rename `config.example.yml` as `config.yml`, then set your own values in `config.yml` according to the following:

#### Trello Client Configuration
You need to set your [Trello API key & token](https://trello.com/app-key) in the configuraiton file, as well as the Trello board & list IDs.

The given list will be the one that new cards is going to be inserted, and it has to be in the given board.

#### Disbling Individual Data Sources
In order to disable a source, just update the `enabled` flag to `false`. There's no need to remove/edit the other parameters for that source.

#### Defining Custom Periods
You can define a custom query period for each source, by populating the `type` and `interval` fields under the `period` for a source.

Example:
```yml
# query at 3rd, 6th, 9th, ... of each month
period:
  type: day
  interval: 3

# query at 00:00, 02:00, 04:00, ... every day
period:
  type: hour
  interval: 2

# query at XX:00, XX:15, XX:30 and XX:45 every hour
period:
  type: minute
  interval: 15

# query on each execution
period:
  type: default
  interval:
```

## Example Cron Job
It's important to make sure that the cron job runs frequently enough to accomodate the most frequent custom interval for a source. It wouldn't make sense to define a custom period of 15 minutes while the cron job only runs every hour.

Both of the following jobs run every hour and both assume that `config.yml` is located in the current working directory.
``` sh
# use "go run"
# 'config.yml' should be located in '/home/utku/git/entrello'
# your go executable may or may not be located in the same place (i.e. /usr/local/go/bin/)
0 * * * * cd /home/utku/git/entrello && /usr/local/go/bin/go run ./cmd/entrello

# use binary executable
# see releases: https://github.com/utkuufuk/entrello/releases
# 'config.yml' should be located in '/path/to/binary'
0 * * * * cd /path/to/binary && ./entrello
```

## 3rd Party Dependencies
| Dependency | Purpose |
|:-|:-|
| [adlio/trello](https://github.com/adlio/trello)           | Trello API Client |
| [google/go-cmp](https://github.com/google/go-cmp)         | Equality Comparisons in Tests |
| [go-github/github](https://github.com/google/go-github)   | GitHub API Client |
| [golang/oauth2](https://github.com/golang/oauth2)         | OAuth 2.0 Client |
| [go-yaml/yaml](https://github.com/go-yaml/yaml)           | Decoding YAML Configuration |
