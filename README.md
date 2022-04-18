# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

## Table of Contents
- [Features](#features)
- [Server Mode Configuration](#server-mode-configuration)
- [Runner Mode Configuration](#runner-mode-configuration)
- [Service Configuration](#service-configuration)
- [Running With Docker](#running-with-docker)
- [Trello Webhooks Reference](#trello-webhooks-reference)

---

## Features
`entrello` synchronizes all your tasks from various sources in one Trello board. It also lets you build automations that can be triggered via the Trello UI.

It can be used either as a **server** or a **runner** (e.g. a cronjob).

### Synchronization
`entrello` synchronizes your tasks from one or more sources in one Trello board:
1. Polls one or more HTTP services of your own, each of which must return a JSON array of "tasks".
2. Creates a new card in your Trello board for each new task it has received from your services. Optionally, it can also remove any stale cards.

Synchronization feature is supported by both [runner](#runner-mode-configuration) and [server](#server-mode-configuration) modes.

### Automation
`entrello` can trigger your HTTP services whenever a card is archived via Trello UI:
1. When a user archives a card via Trello UI, it forwards this event to the matching HTTP service, if any.
2. The matching HTTP service must handle incoming `POST` requests from `entrello` to react to events.

Automation feature is supported only by the [server](#server-mode-configuration) mode because `entrello` needs to expose a callback URL for Trello webhooks.

---

## Server Mode Configuration
Copy and rename `.env.example` as `.env`, then set your own values in `.env`.

You can trigger a poll by making a `POST` request to the root URL of your server with the [service configuration](#service-configuration) in the request body:

```sh
# start the server
go run ./cmd/server

# make a `POST` request to the HTTP server to trigger polling
curl -d @<path/to/config.json> <SERVER_URL> -H "Authorization: Basic <base64(<USERNAME>:<PASSWORD>)>"
```

You can create a [Trello webhook](#trello-webhooks-reference) pointed at `<SERVER_URL>/trello-webhook` in order to listen to events from your Trello board.

---

## Runner Mode Configuration
Create a [service configuration](#service-configuration) file (`config.json` by default), based on `config.example.json`.

You can also use a custom config file path using the `-c` flag:
```sh
go run ./cmd/runner -c /path/to/config/file

# defaults to ./config.json
go run ./cmd/runner
```

---

## Service Configuration
Each service must return a JSON array of Trello card objects (see `pkg/trello/trello.go`) upon a `GET` request. 

Here's a list of open-source HTTP services that are compatible with `entrello`:
- [utkuufuk/github-service](https://github.com/utkuufuk/github-service)

For each service, the following configuration parameters have to be specified:

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

---

## Running With Docker
A new Docker image will be created upon each release.

*See `.github/workflows/release.yml` for continuous delivery workflow configuration.*

1. Login
    ```sh
    echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
    ```

2. Pull the docker image
    ```sh
    docker pull ghcr.io/utkuufuk/entrello/image:latest
    ```

3. Spawn a container:
    ```sh
    # server mode
    docker run -d \
        --env-file <path/to/.env> \
        -p <PORT>:<PORT> \
        --restart unless-stopped \
        --name entrello \
        ghcr.io/utkuufuk/entrello/image:latest

    # runner mode
    docker run --rm \
        -v <path/to/config.json>:/bin/config.json \
        ghcr.io/utkuufuk/entrello/image:latest \
        ./runner
    ```

---

## Trello Webhooks Reference
```sh
# create new webhook
curl -X POST -H "Content-Type: application/json" \
https://api.trello.com/1/tokens/<api_token>/webhooks/ \
-d '{
  "key": "<api_key>",
  "callbackURL": "<callback_url>",
  "idModel": "<board_id>",
  "description": "<desc>"
}'

# list all webhooks
curl https://api.trello.com/1/members/me/tokens?webhooks=true&key=<api_key>&token=<api_token>

# delete existing webhook
curl -X DELETE https://api.trello.com/1/webhooks/<webhook_id>?key=<api_key>&token=<api_token>
```

* `api_token` &mdash; Trello API token
* `api_key` &mdash; Trello API key
* `board_id` &mdash; Trello board ID
* `callback_url` &mdash; Entrello endpoint to handle webhooks, ending with `/trello-webhook`
* `desc` &mdash; Arbitrary description string
* `webhook_id` &mdash; Trello webhook ID

For more information:
* [Trello Webhooks Guide](https://developer.atlassian.com/cloud/trello/guides/rest-api/webhooks/)
* [Trello Webhooks Reference](https://developer.atlassian.com/cloud/trello/rest/#api-group-Webhooks)
