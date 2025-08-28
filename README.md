# Token Bucket

A Redis-based token bucket rate limiter implementation in Go with Lua script for atomic operations.

## Overview

This package provides a distributed token bucket rate limiting mechanism using Redis as the backend store. The token bucket algorithm is ideal for controlling the rate of requests while allowing for burst traffic up to a specified capacity.

## Installation

```bash
go get github.com/rushikeshg25/token-bucket
```

# Token Bucket

A Redis-based token bucket rate limiter implementation in Go with Lua script for atomic operations.

## Overview

This package provides a distributed token bucket rate limiting mechanism using Redis as the backend store. The token bucket algorithm is ideal for controlling the rate of requests while allowing for burst traffic up to a specified capacity.

## Features

- **Distributed rate limiting** using Redis for state persistence
- **Atomic operations** using Lua scripts to ensure consistency
- **Global configuration** with dynamic key-based rate limiting
- **Automatic key expiration** to prevent memory leaks
- **Simple and lightweight** with minimal dependencies
- **Context-aware** operations

## Installation

```bash
go get github.com/rushikeshg25/token-bucket
```

## Usage

### Recommended Usage (Global Rate Limiter)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/redis/go-redis/v9"
    "github.com/rushikeshg25/token-bucket"
)

func main() {
    // Create Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create global rate limiter with configuration
    // Capacity: 100 tokens, Refill rate: 10 tokens per second
    rateLimiter := tokenbucket.NewRateLimiter(rdb, 100, 10)

    ctx := context.Background()

    // Check if request is allowed for specific user
    userID := "user123"
    allowed, err := rateLimiter.Allow(ctx, "user:"+userID)
    if err != nil {
        log.Fatal(err)
    }

    if allowed {
        fmt.Println("Request allowed for user:", userID)
        // Process the request
    } else {
        fmt.Println("Rate limit exceeded for user:", userID)
        // Handle rate limit
    }
}
```

### Middleware Example

```go
import "github.com/gin-gonic/gin"

var RateLimiter *tokenbucket.RateLimiter

func init() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    RateLimiter = tokenbucket.NewRateLimiter(rdb, 1000, 50) // 1000 capacity, 50/sec refill
}

func rateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := getUserIDFromRequest(c) // Your auth logic

        // Each user gets their own rate limit bucket
        allowed, err := globalRateLimiter.Allow(c.Request.Context(), "user:"+userID)
        if err != nil {
            c.JSON(500, gin.H{"error": "Rate limit check failed"})
            c.Abort()
            return
        }

        if !allowed {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Multi-tier Rate Limiting

```go
var (
    apiRateLimiter    *tokenbucket.RateLimiter
    loginRateLimiter  *tokenbucket.RateLimiter
    uploadRateLimiter *tokenbucket.RateLimiter
)

func init() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    apiRateLimiter = tokenbucket.NewRateLimiter(rdb, 1000, 50)    // API calls
    loginRateLimiter = tokenbucket.NewRateLimiter(rdb, 5, 1)      // Login attempts
    uploadRateLimiter = tokenbucket.NewRateLimiter(rdb, 10, 1)    // File uploads
}

func loginHandler(c *gin.Context) {
    username := c.PostForm("username")

    // Rate limit login attempts per user
    allowed, err := loginRateLimiter.Allow(c.Request.Context(), "login:"+username)
    if !allowed || err != nil {
        c.JSON(429, gin.H{"error": "Too many login attempts"})
        return
    }

    // Process login...
}
```

## API Reference

### TokenBucket

#### NewTokenBucket

```go
func NewTokenBucket(redisClient *redis.Client, key string, capacity int, refillRate int) *TokenBucket
```

Creates a new token bucket instance.

**Parameters:**

- `redisClient`: Redis client instance
- `key`: Unique identifier for the bucket (e.g., user ID, IP address)
- `capacity`: Maximum number of tokens the bucket can hold
- `refillRate`: Number of tokens added per second

#### RateLimit

```go
func (tb *TokenBucket) RateLimit(ctx context.Context) (bool, error)
```

Attempts to consume one token from the bucket.

**Returns:**

- `bool`: `true` if request is allowed, `false` if rate limited
- `error`: Error if operation fails

## Algorithm Details

The token bucket algorithm works as follows:

1. **Initialize**: Start with a full bucket of tokens
2. **Refill**: Add tokens at a constant rate (refill rate)
3. **Consume**: Each request consumes one token
4. **Allow/Deny**: Allow request if tokens are available, deny otherwise

### Status Monitoring

```go
func getRateLimitStatus(c *gin.Context) {
    userID := getUserIDFromRequest(c)
    key := "user:" + userID

    tokens, lastRefill, err := globalRateLimiter.GetStatus(c.Request.Context(), key)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to get status"})
        return
    }

    c.JSON(200, gin.H{
        "remaining_tokens": tokens,
        "last_refill_time": lastRefill,
        "capacity":         1000, // Your configured capacity
        "refill_rate":      50,   // Your configured refill rate
    })
}
```
