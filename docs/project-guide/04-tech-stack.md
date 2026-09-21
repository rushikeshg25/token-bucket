# Tech stack

## Languages and runtimes

| Language or runtime | Version | Evidence |
| --- | --- | --- |
| Go | Module directive `1.24.6` | [go.mod:3](../../go.mod#L3) |
| Lua inside Redis | Server-provided, not pinned | Embedded [script](../../lib.go#L19) executed by [EVAL](../../lib.go#L45) |
| Redis server | Not pinned or provisioned here | Caller-supplied [client](../../lib.go#L15) and [deployment limits](../../README.md#L14) |

## Frameworks and major libraries

| Library | Version | Used for | Evidence |
| --- | --- | --- | --- |
| `github.com/redis/go-redis/v9` | `v9.12.1` | Runtime client and integer script-response conversion | [pin](../../go.mod#L11), [call](../../lib.go#L45) |
| `github.com/alicebob/miniredis/v2` | `v2.35.0` | Isolated test Redis and controlled clock | [pin](../../go.mod#L6), [setup](../../lib_test.go#L13) |
| `github.com/yuin/gopher-lua` | `v1.1.1` | Lua dependency in the emulator dependency graph, not imported by package source | [pin](../../go.mod#L13), [test import](../../lib_test.go#L5) |
| Go standard library | Go toolchain | Context/error handling and test concurrency with `sync`/`sync/atomic` | [runtime imports](../../lib.go#L3), [test imports](../../lib_test.go#L3) |

The manifest also lists xxhash, rendezvous hashing, testify, go-spew, go-difflib and yaml as indirect dependencies ([go.mod:7](../../go.mod#L7)). None is directly imported by the two package source files. All requirements carry `// indirect`, including Redis and miniredis despite their direct source imports, so those annotations do not accurately distinguish current direct use.

## Data and infrastructure

Redis holds bucket hashes and expiration, provides the clock, and executes the atomic algorithm ([script](../../lib.go#L20)). There is no local database, file storage, queue, service framework or deployment configuration. Redis connection settings and availability handling belong to the application ([ownership](../../README.md#L14)).

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| Go modules | Resolve package dependencies | [go.mod](../../go.mod#L1) |
| Go test | Behavioral tests using miniredis | [lib_test.go](../../lib_test.go#L21) |
| Make | Convenience target for verbose tests | [Makefile:4](../../Makefile#L4) |
| Go race detector | Check test-time concurrent memory access | [documented command](../../README.md#L16) |

No dedicated linter configuration, formatter command, CI workflow or release/build script appears in the [inventory](03-structure.md). Production Redis compatibility and operational behavior remain outside the emulator-based verification evidence.
