# Redis token bucket v1

```go
client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
limiter := tokenbucket.NewRateLimiter(client, 100, 10)
allowed, err := limiter.Allow(ctx, "v1:api:user123")
tokens, lastRefill, err := limiter.GetStatus(ctx, "v1:api:user123")
```

Capacity is the burst size; refill rate is tokens per second. Both must be integers in 1..1e9. Invalid configuration, nil clients and empty keys return errors from operations. Redis TIME supplies microsecond timestamps, so client clocks cannot create tokens. Fractional credit is preserved across frequent requests. Refill, consume, state write and expiration happen in one Lua script.

`GetStatus` refills without consuming, returning whole available tokens and a Unix-seconds timestamp. It creates state for new keys. State expires after one full refill period plus one second; expiry therefore never resets a bucket before it could naturally become full. Conflicting capacities/rates for the same live key return an error. Use separate key namespaces for separate policies. Use a fresh `v1:` namespace when migrating the old second-resolution representation.

The caller owns the Redis client and context. Redis errors fail the operation; the caller chooses how to handle unavailability. Atomicity depends on one authoritative Redis key; asynchronous failover can lose recently consumed state. V1 does not promise a strict global quota through Redis data loss.

`go test -race ./...` uses an isolated miniredis instance, including real Lua execution, controlled time and simultaneous requests. A production Redis deployment was not exercised during this delivery.
