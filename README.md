# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Run this as a cron job to periodically check custom data sources and automatically create Trello cards based on custom filters.

An example use case could be to create a Trello card for each GitHub issue that's been assigned to you.

## Configuration
Copy and rename `config.example.yml` as `config.yml`, then set your own values in `config.yml`. Most of the configuration parameters are self explanatory, so the following only covers some of them:

### Global Timeout
You can edit the `timeout_secs` config value in order to update global timeout (in seconds) for a single execution.

The execution will not terminate until the timeout is reached, so it's important that the timeout is shorter than the cron job period.

### Trello
You need to set your [Trello API key & token](https://trello.com/app-key) in the configuraiton file, as well as the Trello board & list IDs.

The given list will be the one that new cards is going to be inserted, and it has to be in the given board.

### Telegram
You need a Telegram token & a chat ID in order to enable the integration if you want to receive messages on card updates & possible errors.

### Data Sources
Every data source must have the following configuration parameters under the `source_config` key:
* `name`
* `endpoint`
* `strict`
* `label_id`
* `list_id`
* `period`

#### **`strict`**
Strict mode, which is recommended for most cases, can be enabled for individual data sources by setting the `strict` flag to `true`.

When strict mode is enabled, all the existing Trello cards in the board with the label for the corresponding data source will be deleted, unless the card also exists in the fresh data.

For instance, strict mode can be used to automatically remove resolved GitHub issues from the board. Every time the source is queried, it will return an up-to-date set of open issues. If the board contains any cards that doesn't exist in that set, they will be automatically deleted.

#### **`label_id`**
Each data source must have a distinct Trello label associated with it.

#### **`list_id`**
Each data source must have a target Trello list ID associated with it.

#### **`period`**
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
