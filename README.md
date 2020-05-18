# entrello
![build](https://github.com/utkuufuk/entrello/workflows/entrello/badge.svg?branch=master)
![Latest GitHub release](https://img.shields.io/github/release/utkuufuk/entrello.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/utkuufuk/entrello)](https://goreportcard.com/report/github.com/utkuufuk/entrello)

Run this as a cron job to periodically check custom data sources and automatically create Trello cards based on custom filters.

An example use case (which is already implemented) could be to create a Trello card for each GitHub issue that's been assigned to you.

### Configuration
Copy and rename `config.example.yml` as `config.yml`, then set your own values in `config.yml`.

#### Disbling Individual Data Sources
In order to disable a data source, just update the corresponding line as:
```yml
enabled: false
```
There's no need to edit the remaining config parameters.

### 3rd Party Dependencies
| Dependency | Purpose |
|:-|:-|
| [adlio/trello](https://github.com/adlio/trello)           | Trello API Client |
| [golang/oauth2](https://github.com/golang/oauth2)         | OAuth 2.0 Client |
| [go-github/github](https://github.com/google/go-github)   | GitHub API Client |
| [google/go-cmp](https://github.com/google/go-cmp)         | Equality Comparisons in Tests |

### Example Cron Job
``` sh
# checks data sources every hour
# assumes that `config.yml` is located in `/home/utku/git/entrello`
0 * * * * cd /home/utku/git/entrello && /usr/local/go/bin/go run ./cmd/entrello
```
