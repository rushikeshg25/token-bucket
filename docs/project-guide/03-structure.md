# Structure

## What lives where

The inspected revision has nine tracked files, including two Go source files. All runtime and test code is in the repository root. The working tree also contains a decision log under `logs/`; this guide is the only content added under `docs/project-guide/`.

```text
lib.go                 public API and embedded Lua
lib_test.go            Redis emulator algorithm tests
go.mod / go.sum        module and dependencies
Makefile               verbose test target
README.md              usage and operational limits
V1.md                  delivery contract
HISTORY.md             recorded development history
logs/2026-09-21.md      local decision record
docs/project-guide/    this six-file guide
```

## Root package

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [lib.go](../../lib.go#L1) | Hold policy, validate, run atomic Lua, decode results | `RateLimiter`, `NewRateLimiter`, `Allow`, `GetStatus` | Embedding applications and tests |
| [lib_test.go](../../lib_test.go#L1) | Fixed-clock miniredis setup, fractional refill, concurrent capacity, validation and isolation | `TestFractionalRefill`, `TestConcurrentCapacity`, `TestValidationAndIsolation` | Go test runner |
| [go.mod](../../go.mod#L1) | Declare module, Go version and dependency versions | Module path | Go tooling |
| [Makefile](../../Makefile#L1) | Run `go test -v ./...` | `test` target | Developer invoking make |
| [README.md](../../README.md#L1) | Public example, semantics, deployment caveats | None | Library users |
| [V1.md](../../V1.md#L1) | v1 contract and acceptance | None | Maintainers |
| [HISTORY.md](../../HISTORY.md#L1) | Baseline-to-delivery chronology and evidence | None | Maintainers |

## Logs

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [2026-09-21.md](../../logs/2026-09-21.md#L2) | Local scope and implementation decisions | None | Maintainers |

## Guide

| File | Responsibility | Key exports | Called by |
| --- | --- | --- | --- |
| [README.md](README.md) | Orientation, tour, run commands and open questions | None | New readers |
| [01-architecture.md](01-architecture.md) | Components, state and boundaries | None | Guide readers |
| [02-flow.md](02-flow.md) | Startup, operations and verification | None | Guide readers |
| [03-structure.md](03-structure.md) | File inventory | None | Guide readers |
| [04-tech-stack.md](04-tech-stack.md) | Dependencies and tools | None | Guide readers |
| [05-decisions.md](05-decisions.md) | Evidence, tradeoffs and gotchas | None | Guide readers |

## Excluded

Dependency checksums in [go.sum](../../go.sum) and routine ignore rules in [.gitignore](../../.gitignore) do not define runtime components. There is no vendored source, CI workflow, executable, deployment manifest or generated runtime output in the inspected tree.
