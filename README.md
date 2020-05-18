# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)

Run this as a cron job to periodically check custom data sources and automatically create Trello cards based on custom filters.

An example use case (which is already implemented) could be to create a Trello card for each GitHub issue that's been assigned to you.

### Currently Available Sources
 * Github Issues - https://github.com/issues/assigned
 * TodoDock Tasks - https://tododock.com

Feel free to add new sources or improve the implementations of the existing ones. Contributions are always welcome!

### Configuration
Copy and rename `config.example.yml` as `config.yml`, then set your own values in `config.yml`.

#### Disbling Individual Data Sources
In order to disable a source, just update the `enabled` flag as follows. (There's no need to remove/edit the other parameters for that source.)
```yml
enabled: false
```

### Example Cron Job
Both of the following jobs check data sources every hour. Also, they both assume that `config.yml` is located in the current working directory.
``` sh
# using "go run"
0 * * * * cd /home/utku/git/entrello && /usr/local/go/bin/go run ./cmd/entrello

# using binary executable (see releases: https://github.com/utkuufuk/entrello/releases)
0 * * * * cd /path/to/binary && ./entrello
```

### 3rd Party Dependencies
| Dependency | Purpose |
|:-|:-|
| [adlio/trello](https://github.com/adlio/trello)           | Trello API Client |
| [golang/oauth2](https://github.com/golang/oauth2)         | OAuth 2.0 Client |
| [go-github/github](https://github.com/google/go-github)   | GitHub API Client |
| [google/go-cmp](https://github.com/google/go-cmp)         | Equality Comparisons in Tests |
