# sp0rkle Agent Guide: Instructions for Developers and AI Agents

Welcome to the sp0rkle development guide. This document is designed to help you navigate the codebase, understand its unique architecture, and follow the established coding conventions.

## 1. Architecture Overview

sp0rkle is a modular IRC bot with a driver-based architecture.

- **`bot/`**: The core logic. Manages the IRC connection, command dispatching, rewriters, and pollers.
- **`db/`**: The database abstraction layer. Currently handles the dual-writing logic for the MongoDB to BoltDB migration.
- **`drivers/`**: Modular features (e.g., `factdriver`, `quotedriver`). Each driver is self-contained and registers its functionality with the core bot.
- **`collections/`**: Data access objects for specific features (e.g., `collections/karma`). These sit between the drivers and the `db/` package.
- **`util/`**: General-purpose utilities, including complex lexing (`util/datetime`, `util/lexer.go`).
- **`main.go`**: The entry point. It initializes databases and manually registers drivers.

## 2. Coding Style & "sp0rkle-isms"

This project uses Go 1.22 and follows some non-standard patterns that you must respect:

- **The Global `bot` Singleton**: The `bot` package uses a global singleton to manage state. While not "modern Go," it is the established pattern here. Use `bot.Command`, `bot.Handle`, etc., to register functionality.
- **Manual Registration**: Drivers must be imported in `main.go` and their `Init()` function called explicitly. There is no auto-discovery.
- **Driver Globals**: It is common for drivers to maintain their own package-level state (e.g., collection handles, rate limit maps).
- **Concurrency**: IRC handlers are executed in concurrent goroutines. Shared state in drivers **must** be protected by `sync.Mutex` or `sync.RWMutex`.
- **Go 1.22 Conventions**: Use `interface{}` (not `any`), `math/rand` (seeded in `main.go`), and `io/ioutil` (though `os` is preferred where applicable).
- **Tabs for Indentation**: Use standard `gofmt` with tabs.

## 3. The "Long Slog" Migration (Mongo to BoltDB)

We are in the middle of a migration from MongoDB to BoltDB.

- **`db.Both`**: Most collections use `db.Both`, which writes to both databases simultaneously.
- **`db.K`**: Keys are constructed using a custom key builder to ensure compatibility across both storage engines.
  ```go
  key := db.K{db.S{"nick", "fluffle"}, db.I{"count", 42}}
  ```
- **Dual Writing**: When adding new data-handling logic, ensure it supports the `db.Collection` interface and works with the dual-writing system if applicable.
- **BoltDB as Future**: New features should prioritize BoltDB compatibility.

## 4. Extending the Bot

### Drivers
To add a new feature, create a new package in `drivers/`. It must have an `Init()` function:
```go
func Init() {
    // Register a command: !myfeat <args>
    bot.Command(myHandler, "myfeat", "myfeat <args> -- descriptive help")

    // Register a raw IRC handler
    bot.Handle(rawHandler, client.PRIVMSG)
}
```

### Pollers
If a feature needs to perform periodic background tasks:
1. Implement the `bot.Poller` interface (`Start()`, `Stop()`, `Poll()`, `Tick()`).
2. Register it using `bot.Poll(myPoller)`.
Pollers are automatically started/stopped based on IRC connection status.

### Rewriters
Rewriters allow you to modify outgoing text before it is sent to IRC. Register them with `bot.Rewrite(myRewriter)`.

### Handlers and Context
Handlers receive a `*bot.Context`, which provides:
- `ctx.Text()`: The message body (prefix/name stripped).
- `ctx.ReplyN("format %s", arg)`: Replies to the user with a "Nick: " prefix.
- `ctx.Reply("format %s", arg)`: Replies without the nick prefix.
- `ctx.Line`: The underlying `*client.Line` from `goirc`.

## 5. Testing

- Tests live in `_test.go` files alongside the source.
- **Mocking**: You can often test handler logic by creating a mock `bot.Context` with a manual `client.Line`.
- **Data Tests**: Use the `collections/` package tests as a reference for testing database interactions.

## 6. Tips for AI Agents

- **Trace to Source**: If you see a file in `util/` that looks generated (like `util/datetime/y.go`), look for its source (e.g., `util/datetime/datetime.y`).
- **Check `main.go`**: If your new driver isn't responding, ensure it was added to the `Init` calls in `main.go`.
- **Logging**: Use `github.com/fluffle/golog/logging`.
- **Database Keys**: Be extremely careful with `db.K` composition; changing a key structure can make existing data unreachable.

Happy Hacking!
