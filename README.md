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
- [Open Source `entrello` Services](#open-source-entrello-services)
- [Trello Webhooks Reference](#trello-webhooks-reference)

---

## Features
`entrello` synchronizes all your tasks from various sources in one Trello board. It also lets you build automations that can be triggered via the Trello UI.

It can be used either as a **server** or a **runner** (e.g. a cronjob).

#### Synchronization
`entrello` synchronizes your tasks from one or more sources in one Trello board:
1. Polls one or more HTTP services of your own, each of which must return a JSON array of "tasks".
2. Creates a new card in your Trello board for each new task it has received from your services. Optionally, it can also remove any stale cards.

Synchronization feature is supported by both [runner](#runner-mode-configuration) and [server](#server-mode-configuration) modes.

#### Automation
`entrello` can trigger your HTTP services whenever a card is archived via Trello UI:
1. When a user archives a card via Trello UI, it forwards this event to the matching HTTP service, if any.
2. The matching HTTP service must handle incoming `POST` requests from `entrello` to react to events.

Automation feature is supported only by the [server](#server-mode-configuration) mode, in which a callback URL is exposed for Trello webhooks.

---

## Server Mode Configuration
Put your environment variables in a file called `.env`, based on `.env.example`, and start the server:
```sh
go run ./cmd/server
```

You can trigger a synchronization by making a `POST` request to the root URL of your server with the [service configuration](#service-configuration) in the request body:
```sh
# run this as a scheduled (cron) job
curl -d @<path/to/config.json> <SERVER_URL> -H "Authorization: Basic <base64(<USERNAME>:<PASSWORD>)>"
```

In order to enable the automation feature for one or more services:
1. Create a [Trello webhook](#trello-webhooks-reference), where the callback URL must be `<ENTRELLO_SERVER_URL>/trello-webhook`. 
2. Set the `SERVICES` environment variable to specify a one-on-one mapping of Trello labels to service endpoints.

---

## Runner Mode Configuration
Create a [service configuration](#service-configuration) file based on `config.example.json`. By default, the runner looks for a file called `config.json` in the current working directory.

You can trigger a synchronization by simply executing the runner:
```sh
# run this as a scheduled (cron) job
go run ./cmd/runner
```

Alternatively you can specify a custom config file path using the `-c` flag:
```sh
go run ./cmd/runner -c /path/to/config/file
```

---

## Service Configuration
Each service must return a JSON array of [Trello card objects][1] upon a `GET` request.

For each service, you must set the configuration parameters detailed below:

- `name` &mdash; Service name.

- `endpoint` &mdash; Service endpoint. `entrello` will make `GET` requests to this endpoint to fetch fresh tasks from the service.

- `strict` &mdash; Boolean, whether to delete stale cards or not.

- `label_id` &mdash; Trello label ID associated with the service. A label ID must not be set for more than one service.

- `list_id` &mdash; Trello list ID for the service to determine where to insert new cards. The list must be in the same board specified by the root-level `board_id` config parameter.

- `period` &mdash; Polling period for the service. Determines how often a service should be polled. A few examples:
    ```json
    // poll on 3rd, 6th, 9th, ... of each month, at 00:00
    "period": {
      "type": "day",
      "interval": 3
    }

    // poll every day at 00:00, 02:00, 04:00, ...
    "period": {
      "type": "hour",
      "interval": 2
    }

    // poll every hour at XX:00, XX:15, XX:30, XX:45
    "period": {
      "type": "minute",
      "interval": 15
    }

    // poll on each execution
    "period": {
      "type": "default",
      "interval": 0
    }
    ```

---

## Running With Docker
A new Docker image will be created upon each release.

1. Log in to the GitHub container registry:
    ```sh
    echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
    ```

2. Pull the docker image:
    ```sh
    # find all available versions at
    # https://github.com/utkuufuk/entrello/pkgs/container/entrello%2Fimage/versions
    docker pull ghcr.io/utkuufuk/entrello/image:<tag>

    # or simply download the latest version
    docker pull ghcr.io/utkuufuk/entrello/image:latest
    ```

3. Spawn a container:
    ```sh
    # server mode
    docker run -d \
        --env-file </absolute/path/to/.env> \
        -p <PORT>:<PORT> \
        --restart unless-stopped \
        --name entrello-server \
        ghcr.io/utkuufuk/entrello/image:latest

    # runner mode
    docker run --rm \
        -v </absolute/path/to/config.json>:/bin/config.json \
        ghcr.io/utkuufuk/entrello/image:latest \
        ./runner
    ```

---

## Open Source `entrello` Services
You can use these services directly, or as a reference for developing your own:
- [utkuufuk/github-service](https://github.com/utkuufuk/github-service)
- _stay tuned for more_

---

## Trello Webhooks Reference
```sh
# create new webhook
curl -X POST -H "Content-Type: application/json" -d \
'{
  "key": "<api_key>",
  "callbackURL": "<callback_url>",
  "idModel": "<board_id>",
  "description": "<desc>"
}' https://api.trello.com/1/tokens/<api_token>/webhooks/


# list all webhooks
curl https://api.trello.com/1/members/me/tokens?webhooks=true&key=<api_key>&token=<api_token>

# delete existing webhook
curl -X DELETE https://api.trello.com/1/webhooks/<webhook_id>?key=<api_key>&token=<api_token>
```

| Placeholder   | Description |
|---------------|-------------|
|`api_token`    | Trello API token |
|`api_key`      | Trello API key |
|`board_id`     | Trello board ID |
|`callback_url` | `entrello` server callback URL (see [server config](#server-mode-configuration)) |
|`desc`         | Arbitrary description string |
|`webhook_id`   | Trello webhook ID |

For more information:
* [Trello Webhooks Guide](https://developer.atlassian.com/cloud/trello/guides/rest-api/webhooks/)
* [Trello Webhooks Reference](https://developer.atlassian.com/cloud/trello/rest/#api-group-Webhooks)

[1]: https://github.com/utkuufuk/entrello/blob/master/pkg/trello/trello.go#:~:text=func-,NewCard,-(name%2C%20description%20string
