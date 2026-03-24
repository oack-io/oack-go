# oack-go

Official Go client for the [Oack](https://oack.io) monitoring API.

[![Go Reference](https://pkg.go.dev/badge/github.com/oack-io/oack-go.svg)](https://pkg.go.dev/github.com/oack-io/oack-go)
[![CI](https://github.com/oack-io/oack-go/actions/workflows/ci.yml/badge.svg)](https://github.com/oack-io/oack-go/actions/workflows/ci.yml)

## Install

```sh
go get github.com/oack-io/oack-go
```

Requires Go 1.24+.

## Usage

### API key (Terraform provider, CI/CD)

```go
package main

import (
    "context"
    "fmt"

    oack "github.com/oack-io/oack-go"
)

func main() {
    c := oack.New(oack.BearerToken("sk-your-api-key"))

    teams, err := c.ListAccountTeams(context.Background(), "your-account-id")
    if err != nil {
        panic(err)
    }
    for _, t := range teams {
        fmt.Println(t.ID, t.Name)
    }
}
```

### Dynamic JWT (CLI tools)

```go
c := oack.New(oack.TokenFunc(func() string {
    return currentJWT // refreshed elsewhere
}))

me, err := c.Whoami(context.Background())
```

### Custom base URL

```go
c := oack.New(auth, oack.WithBaseURL("http://localhost:8080"))
```

## Error handling

All non-2xx responses return `*oack.APIError`. Use helpers to check common cases:

```go
monitor, err := c.GetMonitor(ctx, teamID, monitorID)
if oack.IsNotFound(err) {
    // handle 404
}
if oack.IsForbidden(err) {
    // handle 403
}
```

## Resources

| Resource | Methods |
|----------|---------|
| Accounts | Create, List, Get, Update, Delete, Restore, Transfer, Members, Invites, Subscription |
| Teams | Create, List, Get, Update, Delete, Members, Invites, API Keys |
| Monitors | Create, List, Get, Update, Delete, Pause, Unpause, Duplicate, Move, TestAlert |
| Alert Channels | Create, List, Update, Delete, Test, Monitor Links, Alert Events |
| Status Pages | Pages, Components, Groups, Incidents, Maintenance, Subscribers, Templates |
| Watchdogs | Create, List, Update, Delete |
| External Links | Create, List, Get, Update, Delete, Assign/Unassign |
| Integrations | PagerDuty (CRUD + Sync), Cloudflare (CRUD) |
| Probes | List, Get, Details, Aggregate, Download PCAP |
| Comments | Create, List, Edit, Delete, Reply, Resolve, Reopen |
| Notifications | Defaults, Per-Monitor Overrides, Copy Channels |
| Shares | Create, List, Revoke |
| Metrics | Monitor Metrics, Expiration, Timeline, Chart Events |
| Traces | List, Request |
| Geo | Regions, Checkers |
| User | Whoami, Preferences, Devices, Telegram |

## License

[MPL-2.0](LICENSE)
