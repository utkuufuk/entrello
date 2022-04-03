# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

Minimum Go version required: `1.18`

## Usage
- Polls compatible services and keeps your Trello cards synchronized with fresh data.
- Listens for events from your Trello board and forwards "archived card" events to the matching service, if any.
- Can be run as a scheduled job, or as an HTTP server:
    ```sh
    # 1. CRON JOB
    go run ./cmd/runner


    # 2. HTTP SERVER
    go run ./cmd/server

    # make a `POST` request to the HTTP server to trigger polling
    curl -d @config.json <SERVER_URL> -H "Authorization: Basic <base64(<user>:<password>)>"
    ```

### Example Use Case
Let's say you have an HTTP service that returns GitHub issues that are assigned to you upon a `GET` request.
Then `entrello` can use it as a data source to keep your GitHub issues synchronized in your Trello board.

Moreover, when you use `entrello` as a server, it can forward "archived card" events to your GitHub service.
This means that whenever one of your "GitHub" cards is archived, your GitHub service can be notified and take an action of your choosing, e.g. it could close the corresponding GitHub issue. 

## Runner Configuration
Copy and rename `config.example.json` as `config.json` (default), then set your own values in `config.json`.

You can also use a custom config file path using the `-c` flag:
```sh
go run ./cmd/runner -c /path/to/config/file
```

### Configuring Services in Runner Mode
Each configured service must return a JSON array of Trello card objects upon a `GET` request. See `pkg/trello/trello.go` for reference. For each service, the following configuration parameters have to be specified:

- `name` &mdash; Service name.

- `endpoint` &mdash; Service endpoint. `entrello` will make a `GET` request to this endpoint to fetch fresh cards from the service.

- `strict` &mdash; When strict mode is enabled, previously auto-generated cards that are no longer present in the fresh data will be deleted. For instance, with a GitHub service, strict mode can be useful for automatically removing previously auto-generated cards for issues/PRs from the board when the corresponding issues/PRs are closed/merged.

- `label_id` &mdash; **Distinct** Trello label ID associated with the service.

- `list_id` &mdash; Trello list ID for the service to determine where to insert new cards. The selected list must be in the same board as configured by the `board_id` parameter.

- `period` &mdash; Polling period for the service. Some examples:
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

### Example Runner Cron Job
Assuming `config.json` is located in the current working directory:
``` sh
0 * * * * cd /home/you/git/entrello && /usr/local/go/bin/go run ./cmd/runner
```

Make sure that the cron job runs frequently enough to keep up with the most frequent custom interval in your configuration. For instance, it wouldn't make sense to define a custom period of 15 minutes while the cron job only runs every hour.

## Server Configuration
Copy and rename `.env.example` as `.env`, then set your own values in `.env`.

You can trigger the runner by making a `POST` request to the root URL of your server with the runner configuration in the request body:
```sh
curl -d @config.json <SERVER_URL> -H "Authorization: Basic <base64(<USERNAME>:<PASSWORD>)>"
```

You can create a Trello webhook pointed at `<SERVER_URL>/trello-webhook` in order to listen to events from your Trello board.

### Creating Trello Webhooks
You can create a Trello webhook using the following command:

```sh
curl -X POST -H "Content-Type: application/json" \
https://api.trello.com/1/tokens/<api_token>/webhooks/ \
-d '{
  "key": "<api_key>",
  "callbackURL": "<url>",
  "idModel": "<id_model>",
  "description": "<desc>"
}'
```

* `api_token` &mdash; Trello API token
* `api_key` &mdash; Trello API key
* `url` &mdash; Entrello Server URL
* `id_model` &mdash; Trello Board ID
* `desc` &mdash; Arbitrary description string

For more information, see
* [Trello Webhooks Guide](https://developer.atlassian.com/cloud/trello/guides/rest-api/webhooks/)
* [Trello Webhooks Reference](https://developer.atlassian.com/cloud/trello/rest/#api-group-Webhooks)
