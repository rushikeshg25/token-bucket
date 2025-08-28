package tokenbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucket struct {
	redisClient *redis.Client
	key         string
	capacity    int
	refillRate  int
}

func NewTokenBucket(redisClient *redis.Client, key string, capacity int, refillRate int) *TokenBucket {
	return &TokenBucket{
		redisClient: redisClient,
		key:         key,
		capacity:    capacity,
		refillRate:  refillRate,
	}
}

func (tb *TokenBucket) RateLimit(ctx context.Context) (bool, error) {
	luaScript := `
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local refill = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])

	local data = redis.call("HMGET", key, "tokens", "last_refill")
	local tokens = tonumber(data[1])
	local last_refill = tonumber(data[2])

	if tokens == nil then
		tokens = capacity
		last_refill = now
	end

	local delta = now - last_refill
	local refill_tokens = math.floor(delta * refill)
	tokens = math.min(tokens + refill_tokens, capacity)
	last_refill = now 

	local allowed = 0
	if tokens > 0 then
		tokens = tokens - 1
		allowed = 1
	end

	redis.call("HMSET", key, "tokens", tokens, "last_refill", last_refill)
	return allowed
	`
	now := time.Now().Unix()
	res, err := tb.redisClient.Eval(ctx, luaScript, []string{tb.key},
		tb.capacity, tb.refillRate, now).Result()
	if err != nil {
		return false, fmt.Errorf("failed to eval token bucket: %w", err)
	}

	allowed, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected result")
	}

	return allowed == 1, nil
}
