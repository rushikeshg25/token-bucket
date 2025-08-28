package tokenbucket

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket(t *testing.T) {
	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	ctx := context.Background()
	redis.Del(ctx, "test-key")

	tb := NewRateLimiter(redis, 5, 2)

	// First 5 tokens allowed
	for i := 0; i < 5; i++ {
		allowed, err := tb.Allow(ctx, "test-key")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// Next token denied
	allowed, err := tb.Allow(ctx, "test-key")
	assert.NoError(t, err)
	assert.False(t, allowed)

	// 2 tokens refilled
	time.Sleep(time.Second * 1)

	allowed, err = tb.Allow(ctx, "test-key")
	assert.NoError(t, err)
	assert.True(t, allowed) // refill allowed

	allowed, err = tb.Allow(ctx, "test-key")
	assert.NoError(t, err)
	assert.True(t, allowed) // refill allowed

	allowed, err = tb.Allow(ctx, "test-key")
	assert.NoError(t, err)
	assert.False(t, allowed) // no more tokens

}
