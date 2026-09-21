package tokenbucket

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient          *redis.Client
	capacity, refillRate int
}

func NewRateLimiter(client *redis.Client, capacity, refillRate int) *RateLimiter {
	return &RateLimiter{client, capacity, refillRate}
}

const script = `
local capacity, rate, consume = tonumber(ARGV[1]), tonumber(ARGV[2]), tonumber(ARGV[3])
local clock = redis.call('TIME')
local now = tonumber(clock[1]) * 1000000 + tonumber(clock[2])
local state = redis.call('HMGET', KEYS[1], 'tokens', 'last_refill_us', 'capacity', 'rate')
local tokens, last = tonumber(state[1]), tonumber(state[2])
if state[3] and (tonumber(state[3]) ~= capacity or tonumber(state[4]) ~= rate) then
 return redis.error_reply('bucket configuration mismatch')
end
if not tokens or not last then tokens, last = capacity, now end
now = math.max(now, last)
tokens = math.min(capacity, tokens + (now - last) * rate / 1000000)
local allowed = 0
if consume == 1 and tokens >= 1 then tokens = tokens - 1; allowed = 1 end
redis.call('HSET', KEYS[1], 'tokens', tokens, 'last_refill_us', now, 'capacity', capacity, 'rate', rate)
redis.call('PEXPIRE', KEYS[1], math.ceil(capacity / rate * 1000) + 1000)
return {allowed, math.floor(tokens), math.floor(now / 1000000)}
`

func (rl *RateLimiter) eval(ctx context.Context, key string, consume int) ([]int64, error) {
	if rl == nil || rl.redisClient == nil || rl.capacity <= 0 || rl.refillRate <= 0 || rl.capacity > 1000000000 || rl.refillRate > 1000000000 {
		return nil, errors.New("capacity and refill rate must be 1..1e9 and Redis client must be non-nil")
	}
	if key == "" {
		return nil, errors.New("bucket key is empty")
	}
	result, err := rl.redisClient.Eval(ctx, script, []string{key}, rl.capacity, rl.refillRate, consume).Int64Slice()
	if err != nil {
		return nil, fmt.Errorf("token bucket: %w", err)
	}
	if len(result) != 3 {
		return nil, errors.New("unexpected Redis response")
	}
	return result, nil
}
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	r, e := rl.eval(ctx, key, 1)
	if e != nil {
		return false, e
	}
	return r[0] == 1, nil
}

// GetStatus atomically refills without consuming; lastRefill is Unix seconds.
func (rl *RateLimiter) GetStatus(ctx context.Context, key string) (tokens, lastRefill int64, err error) {
	r, e := rl.eval(ctx, key, 0)
	if e != nil {
		return 0, 0, e
	}
	return r[1], r[2], nil
}
