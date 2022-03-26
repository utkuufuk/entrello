# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Polls compatible data sources and keeps your Trello cards synchronized with fresh data. Meant to be run as a scheduled job.

Let's say you have an HTTP endpoint that returns GitHub issues assigned to you upon a `GET` request.
You can point `entrello` to that endpoint to keep your GitHub issues synchronized in your Trello board.

Each data source must return a JSON array of Trello card objects upon a `GET` request. You can import and use the `NewCard` function from `pkg/trello/trello.go` in order to construct Trello card objects.

## Configuration
Copy and rename `config.example.json` as `config.json` (default), then set your own values in `config.json`.

You can also use a custom config file path using the `-c` flag:
```sh
go run ./cmd/entrello -c /path/to/config/file
```

Alternatively, you can store the configuration inside the `ENTRELLO_CONFIG` environment variable as a JSON string.

### Parameters
#### Trello
You need to set your [Trello API key & token](https://trello.com/app-key) in the configuraiton file, as well as the Trello board ID.

#### Data Sources
Each data source must have the following configuration parameters. (See `config.example.json`)

#### `name`
Data source name.

#### `endpoint`
Data source endpoint. `entrello` will make a `GET` request to this endpoint to fetch fresh cards from the data source. 

#### `strict`
When strict mode is enabled, previously auto-generated cards that are no longer present in the fresh data will be deleted.

For instance, with a GitHub data source, strict mode can be useful for automatically removing previously auto-generated cards for issues/PRs from the board when the corresponding issues/PRs are closed/merged.

#### `label_id`
**Distinct** Trello label ID associated with the data source.

#### `list_id`
Trello list ID for the data source to determine where to insert new cards. The selected list must be in the same board as configured by the `board_id` parameter.

#### `period`
Polling period for the data source.

Example periods:
```json
// query at 3rd, 6th, 9th, ... of each month
"period": {
  "type": "day",
  "interval": 3
}

// query at 00:00, 02:00, 04:00, ... every day
"period": {
  "type": "hour",
  "interval": 2
}

// query at XX:00, XX:15, XX:30 and XX:45 every hour
"period": {
  "type": "minute",
  "interval": 15
}

// query on each execution
"period": {
  "type": "default",
  "interval": 0
}
```

## Example Cron Job
Make sure that the cron job runs frequently enough to keep up with the most frequent custom interval in your configuration.

For instance, it wouldn't make sense to define a custom period of 15 minutes while the cron job only runs every hour.

Both of the following jobs run every hour and both assume that `config.json` is located in the current working directory, or its contents are stored within the `ENTRELLO_CONFIG` environment variable.
``` sh
# using "go run"
0 * * * * cd /home/you/git/entrello && /usr/local/go/bin/go run ./cmd/entrello

# use binary executable (see releases: https://github.com/you/entrello/releases)
0 * * * * cd /path/to/binary && ./entrello
```
