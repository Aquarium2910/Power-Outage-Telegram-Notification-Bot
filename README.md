# Power Outage Telegram Notification Bot

Telegram bot that notifies users when power outage schedules are updated.
This is a rewrite focused on practicing Clean Architecture in Go — the
codebase is structured into distinct layers with explicit interfaces
between them.

**Status:** early stage — core layers and DB schema are in place;
bot currently only responds to `/start`. Actively in progress.

## Architecture

```text
internal/
├── entity/       # domain models (user, schedule)
├── usecase/      # business logic, defined via interfaces (contracts.go)
├── repo/         # data access, persistent implementation + contracts.go
└── controller/   # delivery layer (telegram handlers, router, presenter)
```

Each layer depends only on interfaces defined by the layer below it
(contracts.go files), not on concrete implementations.

## Stack

Go, PostgreSQL, Docker Compose