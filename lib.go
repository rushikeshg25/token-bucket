package tokenbucket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	capacity    int
	refillRate  int
}

func NewRateLimiter(redisClient *redis.Client, capacity int, refillRate int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		capacity:    capacity,
		refillRate:  refillRate,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
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
	redis.call("EXPIRE", key, 3600) -- Expire after 1 hour of inactivity
	return allowed
	`

	now := time.Now().Unix()
	res, err := rl.redisClient.Eval(ctx, luaScript, []string{key},
		rl.capacity, rl.refillRate, now).Result()
	if err != nil {
		return false, fmt.Errorf("failed to eval rate limiter: %w", err)
	}

	allowed, ok := res.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected result from redis")
	}

	return allowed == 1, nil
}

func (rl *RateLimiter) GetStatus(ctx context.Context, key string) (tokens int64, lastRefill int64, err error) {
	result, err := rl.redisClient.HMGet(ctx, key, "tokens", "last_refill").Result()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get status: %w", err)
	}

	if result[0] != nil {
		tokensStr, ok := result[0].(string)
		if ok {
			tokensInt, err := strconv.Atoi(tokensStr)
			if err == nil {
				tokens = int64(tokensInt)
			}
		}
	}

	if tokens == 0 && result[0] == nil {
		tokens = int64(rl.capacity)
	}

	if result[1] != nil {
		lastRefillStr, ok := result[1].(string)
		if ok {
			lastRefillInt, err := strconv.Atoi(lastRefillStr)
			if err == nil {
				lastRefill = int64(lastRefillInt)
			}
		}
	}

	return tokens, lastRefill, nil
}
