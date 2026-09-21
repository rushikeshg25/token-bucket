# Token bucket project guide

> Generated: 2026-09-21 from commit `9942356`.

## What this is

This Go library limits operations with a Redis-backed token bucket whose capacity sets the burst size and whose rate sets tokens per second ([API](../../lib.go#L10)). Applications supply a Redis client, context and bucket key, then receive an admission decision or an error ([Allow](../../lib.go#L54)). One Lua script uses Redis time to refill, optionally consume, and expire each bucket atomically ([script](../../lib.go#L19)).

## Run it

Use Go 1.24.6 or a compatible newer toolchain ([manifest](../../go.mod#L3)):

```bash
go mod download
go test ./...
go test -race ./...
# Equivalent verbose test target:
make test
```

There is no executable entry point: import `github.com/rushikeshg25/token-bucket` into a Go application and follow the [usage example](../../README.md#L3). Runtime operations require a caller-configured Redis connection with permission for the script's commands; tests start miniredis automatically ([setup](../../lib_test.go#L13)). The library defines no environment variables or Redis deployment configuration.

## The five-file tour

Only two Go source files exist, so the tour continues through the manifest, runner and contract rather than inventing more components.

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [lib.go](../../lib.go#L15) | Constructor, public operations and complete Redis algorithm | [Runtime flow](02-flow.md#allow) |
| 2 | [lib_test.go](../../lib_test.go#L13) | Controlled time and concurrent callers make the guarantees concrete | [Verification flow](02-flow.md#verification) |
| 3 | [go.mod](../../go.mod#L1) | Toolchain and Redis client/emulator versions | [Stack](04-tech-stack.md) |
| 4 | [Makefile](../../Makefile#L1) | The sole checked-in test command | [Tooling](04-tech-stack.md#tooling) |
| 5 | [V1.md](../../V1.md#L3) | Intended v1 boundaries and acceptance criteria | [Decisions](05-decisions.md) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — ownership and Redis state.
2. [Flow](02-flow.md) — construction, admission, status and verification.
3. [Structure](03-structure.md) — complete source inventory.
4. [Tech stack](04-tech-stack.md) — versions and tooling.
5. [Decisions](05-decisions.md) — confirmed choices, inferences and traps.

## Open questions

- Which Redis version, persistence, failover and access-control configuration will host this library? None is pinned here, and [delivery notes](../../README.md#L16) explicitly exclude production Redis validation.
- Should lifecycle tests exercise actual expiry, backward server time, Redis outages and upper-bound inputs? Current tests cover fractional refill, simultaneous capacity, invalid configurations, isolation, conflict, positive TTL and cancellation ([tests](../../lib_test.go#L21)), but do not execute those additional cases.
