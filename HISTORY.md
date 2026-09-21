# Redis token-bucket limiter history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2025-08-28T12:21:19+05:30: init:mod

- **What happened:** The repository records `init:mod`.
- **Evidence:** [commit c991e5b8bb](https://github.com/rushikeshg25/token-bucket/commit/c991e5b8bbacf357b8ef5144cf3e9d178abe4198).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:23:54+05:30: docs: define token-bucket v1 contract

- **What happened:** The repository records `docs: define token-bucket v1 contract`.
- **Evidence:** [commit e584f09e54](https://github.com/rushikeshg25/token-bucket/commit/e584f09e546044cdabcfb0e1e91bf03ca0ea15d6).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:30:39+05:30: fix: use Redis time and retain fractional refill credit atomically

- **What happened:** The repository records `fix: use Redis time and retain fractional refill credit atomically`.
- **Evidence:** [commit 90d8000d14](https://github.com/rushikeshg25/token-bucket/commit/90d8000d143e3a1a7f7aee4b385484f7aedac7e4).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:33:02+05:30: test: verify atomic capacity fractional refill and configuration using Redis emulator

- **What happened:** The repository records `test: verify atomic capacity fractional refill and configuration using Redis emulator`.
- **Evidence:** [commit 570e877a12](https://github.com/rushikeshg25/token-bucket/commit/570e877a12e8987c0378300fe8fadc62717f255b).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:38:57+05:30: docs: document token-bucket v1 usage and limitations

- **What happened:** The repository records `docs: document token-bucket v1 usage and limitations`.
- **Evidence:** [commit 45393448ae](https://github.com/rushikeshg25/token-bucket/commit/45393448ae4141e58e17b3a7ee6a31edc128fbd9).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed `go test -race ./...` using miniredis Lua execution and controlled server time. [lib_test.go](lib_test.go) checks fractional refill, shared capacity, configuration, expiry and cancellation.

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.
