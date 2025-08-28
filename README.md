# Token Bucket

A Redis-based token bucket rate limiter implementation in Go with Lua script for atomic operations.

## Overview

This package provides a distributed token bucket rate limiting mechanism using Redis as the backend store. The token bucket algorithm is ideal for controlling the rate of requests while allowing for burst traffic up to a specified capacity.

## Installation

```bash
go get github.com/rushikeshg25/token-bucket
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/rushikeshg25/token-bucket"
)

func main() {
    // Create Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create token bucket
    // Capacity: 10 tokens, Refill rate: 1 token per second
    tb := tokenbucket.NewTokenBucket(rdb, "user:123", 10, 1)

    ctx := context.Background()

    // Check if request is allowed
    allowed, err := tb.RateLimit(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if allowed {
        fmt.Println("Request allowed")
        // Process the request
    } else {
        fmt.Println("Rate limit exceeded")
        // Handle rate limit
    }
}
```

### API Rate Limiting Example

```go
func rateLimitMiddleware(tb *tokenbucket.TokenBucket) gin.HandlerFunc {
    return func(c *gin.Context) {
        allowed, err := tb.RateLimit(c.Request.Context())
        if err != nil {
            c.JSON(500, gin.H{"error": "Internal server error"})
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

### Key Properties

- **Burst handling**: Allows up to `capacity` requests in quick succession
- **Sustained rate**: Long-term rate is limited to `refillRate` requests per second
- **Distributed**: State is shared across multiple instances via Redis
