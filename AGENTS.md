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

## 3. Events and Handlers

sp0rkle uses an event-based system built on top of `goirc`.

- **Handler Signature**: `type HandlerFunc func(*bot.Context)`
- **`bot.Handle(fn, events...)`**: Registers a `HandlerFunc` for specific IRC events (e.g., `client.PRIVMSG`, `client.JOIN`).
- **`bot.HandleBG(fn, events...)`**: Same as `Handle`, but runs the handler in its own goroutine. Useful for long-running tasks like migrations.
- **`bot.Context`**: The primary object passed to handlers. It encapsulates the IRC line, the connection, and provides helper methods:
    - `ctx.Text()`: Message body with bot name/command prefix stripped.
    - `ctx.ReplyN(format, args...)`: Reply with "Nick: " prefix.
    - `ctx.Storable()`: Returns sender nick and channel.

## 4. Plugins

The Factoid driver supports "plugins" which allow other drivers to perform transformations on factoid values.
- **Implementation**: A driver provides a `RegisterPlugins` method (if using the older pattern) or simply registers functions that the factoid driver calls.
- **Factoid Syntax**: Triggered via `<plugin=name args>` in a factoid value.
- **Identifier Replacement**: Common identifiers like `$nick`, `$chan`, `$date`, and `$time` are handled by a standard replacer in `factdriver/plugins.go`.

## 5. The "Long Slog" Migration (Mongo to BoltDB)

We are in the middle of a migration from MongoDB to BoltDB.

- **`db.Both`**: Most collections use `db.Both`, which writes to both databases simultaneously and compares reads.
- **`db.K`**: Keys are constructed using a custom key builder:
  ```go
  key := db.K{db.S{"nick", "fluffle"}, db.I{"count", 42}}
  ```
- **BoltDB Structure**: Successive key elements create nested BoltDB buckets, with the final element as the key.

## 6. Extending the Bot

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
For periodic tasks, implement `bot.Poller` (`Start`, `Stop`, `Poll`, `Tick`) and register with `bot.Poll(myPoller)`.

## 7. Testing

- Tests live in `_test.go` files.
- Use `bot.Context` mocking for handler tests.
- Reference `collections/` tests for database interaction testing.

## 8. Tips for AI Agents

- **Trace to Source**: Many files in `util/` (like `datetime/y.go`) are generated from `.y` or `.rl` files.
- **Check `main.go`**: Manual registration is required for all drivers.
- **Logging**: Use `github.com/fluffle/golog/logging`.

Happy Hacking!
