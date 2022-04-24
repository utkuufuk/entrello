# entrello
[![CI](https://github.com/utkuufuk/entrello/actions/workflows/ci.yml/badge.svg)](https://github.com/utkuufuk/entrello/actions/workflows/ci.yml)
![Latest Release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Coverage Status](https://coveralls.io/repos/github/utkuufuk/entrello/badge.svg)](https://coveralls.io/github/utkuufuk/entrello)

## Table of Contents
- [Features](#features)
- [Runner Mode Configuration](#runner-mode-configuration)
- [Server Mode Configuration](#server-mode-configuration)
- [Service Configuration](#service-configuration)
- [Running With Docker](#running-with-docker)
- [Available `entrello` Services](#available-entrello-services)
- [Trello Webhooks Reference](#trello-webhooks-reference)

---

## Features
- Synchronizes all your tasks from various sources in one Trello board.
- Lets you build automations that can be triggered by user actions via Trello UI.
- Can be used either as a **server** or a **runner** (e.g. a cronjob).

#### Synchronization
`entrello` synchronizes your tasks from one or more sources in one Trello board:
1. Polls one or more HTTP services of your own, each of which must return a JSON array of "tasks".
2. Creates a new card in your Trello board for each new task it has received from your services. Optionally, it can also remove any stale cards.

Synchronization feature is supported by both the [runner](#runner-mode-configuration) and [server](#server-mode-configuration) modes.

#### Automation
`entrello` can trigger your HTTP services whenever a card is archived via Trello UI:
1. When a user archives a card via Trello UI, it forwards this event to the matching HTTP service, if any.
2. The matching HTTP service must handle incoming `POST` requests from `entrello` to react to events.

Automation feature is supported only by the [server](#server-mode-configuration) mode, in which a callback URL is exposed for Trello webhooks.

---

## Runner Mode Configuration
Create a [service configuration](#service-configuration) file based on `config.example.json`. You can trigger a synchronization by simply executing the runner:
```sh
# run this as a scheduled (cron) job
go run ./cmd/runner -c /path/to/config/file
```

If the `-c` flag is omitted, the runner looks for a file called `config.json` in the current working directory:
```sh
# these two are equivalent:
go run ./cmd/runner
go run ./cmd/runner -c ./config.json
```

---

## Server Mode Configuration
Put your environment variables in a file called `.env`, based on `.env.example`, and start the server:
```sh
go run ./cmd/server
```

You can trigger a synchronization by making a `POST` request to the root URL of your server with the [service configuration](#service-configuration) in the request body:
```sh
# run this as a scheduled (cron) job
curl <SERVER_URL> \
    -d @<path/to/config.json> \
    -H "Authorization: Basic <base64(<USERNAME>:<PASSWORD>)>"
```

### Automation
To enable automation for one or more services:
1. Create a [Trello webhook](#trello-webhooks-reference), where the callback URL is `<ENTRELLO_SERVER_URL>/trello-webhook`.
2. Set the `SERVICES` environment variable, a comma-separated list of service configuration strings:
    * A service configuration string must contain the Trello label ID and the service endpoint:
        ```sh
        # trello label ID: 1234
        # service enpoint URL: http://localhost:3333/entrello
        1234@http://localhost:3333/entrello
        ```
    * It may additionally contain an API secret &ndash; _alphanumeric only_ &ndash; for authentication purposes:
        ```sh
        # the HTTP header "X-Api-Key" will be set to "SuPerSecRetPassW0rd" in each request
        1234:SuPerSecRetPassW0rd@http://localhost:3333/entrello
        ```

---

## Service Configuration
Each service must return a JSON array of [Trello card objects][1] upon a `GET` request.

#### Mandatory configuration parameters

- `name` &mdash; Service name.

- `endpoint` &mdash; Service endpoint URL.

- `label_id` &mdash; Trello label ID. A label ID can be associated with no more than one service.

- `list_id` &mdash; Trello list ID, i.e. where to insert new cards. The list must be in the board specified by the root-level `board_id` config parameter.

- `period` &mdash; Polling period. A few examples:
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
      "type": "default"
    }
    ```

#### Optional configuration parameters

- `secret` &mdash; Alphanumeric API secret. If present, `entrello` will put it in the `X-Api-Key` HTTP header.

- `strict` &mdash; Whether stale cards should be deleted from the board upon synchronization. `false` by default.

---

## Running With Docker
A new [Docker image](https://github.com/utkuufuk?tab=packages&repo_name=entrello) will be created upon each [release](https://github.com/utkuufuk/entrello/releases).

1. Authenticate with the GitHub container registry (only once):
    ```sh
    echo $GITHUB_ACCESS_TOKEN | docker login ghcr.io -u GITHUB_USERNAME --password-stdin
    ```

2. Pull the latest Docker image:
    ```sh
    docker pull ghcr.io/utkuufuk/entrello/image:latest
    ```

3. Spawn & run a container:
    ```sh
    # server
    docker run -d \
        -p <PORT>:<PORT> \
        --env-file </absolute/path/to/.env> \
        --restart unless-stopped \
        --name entrello-server \
        ghcr.io/utkuufuk/entrello/image:latest

    # runner
    docker run --rm \
        -v </absolute/path/to/config.json>:/bin/config.json \
        ghcr.io/utkuufuk/entrello/image:latest \
        ./runner
    ```

---

## Available `entrello` Services
You can use these services directly, or as a reference for developing your own:
- [utkuufuk/github-service](https://github.com/utkuufuk/github-service)

_Stay tuned for more..._

---

## Trello Webhooks Reference
```sh
# create new webhook
curl -X POST -H "Content-Type: application/json" -d \
'{
  "key": "<TRELLO_API_KEY>",
  "callbackURL": "<ENTRELLO_SERVER_CALLBACK_URL>",
  "idModel": "<TRELLO_BOARD_ID>",
  "description": "<DESCRIPTION>"
}' https://api.trello.com/1/tokens/<TRELLO_API_TOKEN>/webhooks/


# list all webhooks
curl https://api.trello.com/1/members/me/tokens?webhooks=true&key=<TRELLO_API_KEY>&token=<TRELLO_API_TOKEN>

# delete existing webhook
curl -X DELETE https://api.trello.com/1/webhooks/<TRELLO_WEBHOOK_ID>?key=<TRELLO_API_KEY>&token=<TRELLO_API_TOKEN>
```

For more information on Trello webhooks:
* [Trello Webhooks Guide](https://developer.atlassian.com/cloud/trello/guides/rest-api/webhooks/)
* [Trello Webhooks Reference](https://developer.atlassian.com/cloud/trello/rest/#api-group-Webhooks)

[1]: https://github.com/utkuufuk/entrello/blob/master/pkg/trello/trello.go#:~:text=func-,NewCard,-(name%2C%20description%20string
