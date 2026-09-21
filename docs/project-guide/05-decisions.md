# Decisions

The choices below separate documented intent from implications inferred from the implementation.

## One script and one clock authority

- **What:** Redis Lua performs refill and consumption using Redis `TIME` in a single invocation.
- **Evidence:** [Script](../../lib.go#L20), [EVAL](../../lib.go#L45), and the [v1 contract](../../V1.md#L5).
- **Why:** Atomic admission and fractional accounting are explicit v1 requirements. Redis time avoids disagreement between client clocks ([README](../../README.md#L10)).
- **Tradeoff:** Concurrent callers share one authority, but every admission depends on Redis and script execution. There is no local fallback.
- **Confidence:** Confirmed by contract and implementation.

## Keep fractional state but return whole tokens

- **What:** The hash retains fractional token credit and microsecond timestamps, while replies floor tokens and convert timestamps to seconds.
- **Evidence:** [Refill and reply](../../lib.go#L28), [status comment](../../lib.go#L62), and [fractional-refill test](../../lib_test.go#L21).
- **Why:** Frequent calls must not discard partial refill credit ([README](../../README.md#L10)). Whole-token status reflects what can currently be consumed.
- **Tradeoff:** Status hides fractional progress and subsecond timing even though internal state preserves them.
- **Confidence:** Fractional accounting is confirmed. The explanation for exposing only whole-token availability is inferred.

## Expire after a full refill period plus slack

- **What:** Each successful script call resets TTL to the full refill period rounded up to milliseconds, plus one second.
- **Evidence:** [PEXPIRE](../../lib.go#L34), [documented rationale](../../README.md#L12).
- **Why:** Expiration should not recreate a full bucket earlier than natural refill would permit.
- **Tradeoff:** Inactive keys are reclaimed safely, while frequent denials or status calls keep state alive. Redis loss can still reset state early.
- **Confidence:** Confirmed by documentation and implementation.

## Bind policy to each live key

- **What:** Persist capacity/rate with the balance and reject a caller using a different policy for that key.
- **Evidence:** [Conflict check](../../lib.go#L25), [test](../../lib_test.go#L77), [namespace guidance](../../README.md#L12).
- **Why:** Inferred: prevent incompatible callers silently rewriting the meaning of shared state.
- **Tradeoff:** Policy changes need a different namespace, deliberate state handling or expiration. Construction itself cannot detect a conflict.
- **Confidence:** Behavior confirmed, rationale inferred.

## Caller owns connection and failure policy

- **What:** Accept a prebuilt `*redis.Client` and propagate operation errors rather than choosing availability behavior for an application.
- **Evidence:** [Constructor](../../lib.go#L15), [error propagation](../../lib.go#L46), [ownership documentation](../../README.md#L14).
- **Why:** The README explicitly assigns Redis unavailability handling to callers.
- **Tradeoff:** Integration stays small, but applications must configure connection lifecycle, contexts and error handling themselves.
- **Confidence:** Confirmed by documentation.

## Gotchas

- **Constructor success does not validate configuration.** Invalid values fail only when an operation reaches [eval](../../lib.go#L38).
- **Status is a write.** It creates missing buckets, refills them and renews TTL through the [same script](../../lib.go#L63).
- **`lastRefill` is in seconds.** The stored field is microseconds, and repeated operations update it even when no token is consumed ([reply](../../lib.go#L33)).
- **A denied request can still change state.** Refill, timestamp and TTL are written after the consume decision ([script](../../lib.go#L30)).
- **Namespaces are the caller's responsibility.** Use fresh v1 keys when migrating old second-resolution state, as [README.md:12](../../README.md#L12) directs. The [hash reader](../../lib.go#L23) contains no migration procedure.
- **Positive TTL is not a full expiry test.** The test checks `TTL > 0`, not expiry timing or reinitialization after expiry ([assertion](../../lib_test.go#L80)).
- **Atomicity does not guarantee a global quota after data loss.** The [documented failover limit](../../README.md#L14) applies even when all clients use this API correctly.

## Conventions

Keep runtime policy changes in the single [implementation file](../../lib.go), and use fixed server time plus isolated per-test clients for algorithm checks ([test helper](../../lib_test.go#L13)). Current tests exercise public API behavior rather than exposing the unexported evaluator. Record uncovered production assumptions in the [open questions](README.md#open-questions), rather than treating emulator success as deployment evidence.
