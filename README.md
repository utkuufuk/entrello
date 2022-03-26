# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Polls compatible data sources and keeps your Trello cards synchronized with fresh data.

Let's say you have an HTTP endpoint that returns GitHub issues assigned to you upon a `GET` request.
You can point `entrello` to that endpoint to keep your GitHub issues synchronized in your Trello board.

Each data source must return a JSON array of Trello card objects upon a `GET` request. You can import and use the `NewCard` function from `pkg/trello/trello.go` in order to construct Trello card objects.

- Can be run as a scheduled job:
    ```sh
    go run ./cmd/runner
    ```
- Can be run as an HTTP server:
    ```sh
    PORT=<port> USERNAME=<user> PASSWORD=<password> go run ./cmd/server
    ```

    In this case, the runner can be triggered by a `POST` request to the server like this:
    ```sh
    curl -d @config.json <SERVER_URL> -H "Authorization: Basic <base64(<user>:<password>)>"
    ```

## Configuration
Copy and rename `config.example.json` as `config.json` (default), then set your own values in `config.json`.

You can also use a custom config file path using the `-c` flag:
```sh
go run ./cmd/runner -c /path/to/config/file
```

### Trello
You need to set your [Trello API key & token](https://trello.com/app-key) in the configuraiton file, as well as the Trello board ID.

### Data Sources
For each data source, the following parameters have to be specified. (See `config.example.json`)

- `name` &mdash; Data source name.

- `endpoint` &mdash; Data source endpoint. `entrello` will make a `GET` request to this endpoint to fetch fresh cards from the data source.

- `strict` &mdash; When strict mode is enabled, previously auto-generated cards that are no longer present in the fresh data will be deleted. For instance, with a GitHub data source, strict mode can be useful for automatically removing previously auto-generated cards for issues/PRs from the board when the corresponding issues/PRs are closed/merged.

- `label_id` &mdash; **Distinct** Trello label ID associated with the data source.

- `list_id` &mdash; Trello list ID for the data source to determine where to insert new cards. The selected list must be in the same board as configured by the `board_id` parameter.

- `period` &mdash; Polling period for the data source. Some examples:
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
Assuming `config.json` is located in the current working directory:
``` sh
0 * * * * cd /home/you/git/entrello && /usr/local/go/bin/go run ./cmd/runner
```

Make sure that the cron job runs frequently enough to keep up with the most frequent custom interval in your configuration. For instance, it wouldn't make sense to define a custom period of 15 minutes while the cron job only runs every hour.
