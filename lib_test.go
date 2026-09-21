package tokenbucket

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func setup(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	m := miniredis.RunT(t)
	m.SetTime(time.Unix(1700000000, 0))
	c := redis.NewClient(&redis.Options{Addr: m.Addr()})
	t.Cleanup(func() { c.Close() })
	return m, c
}
func TestFractionalRefill(t *testing.T) {
	m, c := setup(t)
	rl := NewRateLimiter(c, 1, 2)
	ctx := context.Background()
	if ok, e := rl.Allow(ctx, "bucket"); e != nil || !ok {
		t.Fatal(ok, e)
	}
	for i := 1; i <= 5; i++ {
		m.SetTime(time.Unix(1700000000, int64(i)*100000000))
		ok, e := rl.Allow(ctx, "bucket")
		if e != nil || ok != (i == 5) {
			t.Fatalf("step %d allowed %v err %v", i, ok, e)
		}
	}
	m.SetTime(time.Unix(1700000001, 0))
	n, _, e := rl.GetStatus(ctx, "bucket")
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestConcurrentCapacity(t *testing.T) {
	_, c := setup(t)
	rl := NewRateLimiter(c, 10, 1)
	var wg sync.WaitGroup
	var allowed, failures atomic.Int32
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, e := rl.Allow(context.Background(), "shared")
			if e != nil {
				failures.Add(1)
			}
			if ok {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if allowed.Load() != 10 || failures.Load() != 0 {
		t.Fatal(allowed.Load(), failures.Load())
	}
}
func TestValidationAndIsolation(t *testing.T) {
	m, c := setup(t)
	ctx := context.Background()
	for _, rl := range []*RateLimiter{NewRateLimiter(nil, 1, 1), NewRateLimiter(c, 0, 1), NewRateLimiter(c, 1, -1)} {
		if _, e := rl.Allow(ctx, "k"); e == nil {
			t.Fatal("invalid config accepted")
		}
	}
	rl := NewRateLimiter(c, 1, 1)
	rl.Allow(ctx, "a")
	if ok, e := rl.Allow(ctx, "b"); !ok || e != nil {
		t.Fatal("keys share state")
	}
	if _, e := NewRateLimiter(c, 2, 1).Allow(ctx, "a"); e == nil {
		t.Fatal("config conflict")
	}
	if m.TTL("a") <= 0 {
		t.Fatal("missing expiration")
	}
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := rl.Allow(ctx, "c"); e == nil {
		t.Fatal("cancel ignored")
	}
}
