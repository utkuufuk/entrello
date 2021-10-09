# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Run this as a cron job to periodically poll your data sources via HTTP and automatically create/update/delete Trello cards based on fresh data.

An example use case could be to create a Trello card for each GitHub issue that's assigned to you.

Your data sources must return a JSON array of Trello card objects upon a `GET` request. You can import and use the `NewCard` function from `pkg/trello/trello.go` in order to construct Trello card objects.

## Configuration
Copy and rename `config.example.yml` as `config.yml` (default), then set your own values in `config.yml`.

You can also use a custom config file path using the `-c` flag:
```sh
go run ./cmd/entrello -c /path/to/config/file
```

### Trello
You need to set your [Trello API key & token](https://trello.com/app-key) in the configuraiton file, as well as the Trello board ID.

### Telegram
You need a Telegram token & a chat ID in order to receive alerts in case of errors.

### Data Sources
Each data source must have the following configuration parameters. Refer to `config.example.yml` for examples.

#### **`name`**
Data source name.

#### **`endpoint`**
Data source endpoint. `entrello` will make a `GET` request to this endpoint to fetch fresh cards from the data source. 

#### **`strict`**
When strict mode is enabled, previously auto-generated cards that are no longer present in the fresh data will be deleted.

For instance, with a GitHub data source, strict mode can be useful for automatically removing previously auto-generated cards for issues/PRs from the board when the corresponding issues/PRs are closed/merged.

#### **`label_id`**
**Distinct** Trello label ID associated with the data source.

#### **`list_id`**
Trello list ID for the data source to determine where to insert new cards. The selected list must be in the same board as configured by the `board_id` parameter.

#### **`period`**
Polling period for the data source.

Example configuration:
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
# your go executable may or may not be located in the same location (i.e. /usr/local/go/bin/)
0 * * * * cd /home/utku/git/entrello && /usr/local/go/bin/go run ./cmd/entrello

# use binary executable
# see releases: https://github.com/utkuufuk/entrello/releases
# 'config.yml' should be located in '/path/to/binary'
0 * * * * cd /path/to/binary && ./entrello
```
