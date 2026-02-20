package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisHealthy = true
var lastErrorTime time.Time
var slidingWindowLua = redis.NewScript(`
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window_start = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])
	local duration = tonumber(ARGV[4])

	--1 Remove old hits 
	redis.call("ZREMRANGEBYSCORE", key, 0, window_start)

	--2 Count current hits
	local count = redis.call("ZCARD", key)

	--3. If under limit, add the new hit
	if count< limit then
		redis.call("ZADD", key, now, now)
		redis.call("EXPIRE", key, duration)
		return 0 --success
	else 
		return 1 -- Limit Reached
	end
`)
var ctx = context.Background()
var rdb *redis.Client

func rateLimitMiddleware(next http.HandlerFunc, limit int64, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			return
		}
		userIp, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			userIp = r.RemoteAddr
		}
		limitKey := fmt.Sprintf("%s:%s", userIp, r.URL.Path)
		// Get current time in Microseconds (high precision)
		now := time.Now().UnixMicro()
		windowStartTime := now - duration.Microseconds()
		// 1. If the "Circuit" is open, skip Redis entirely to save time
		if !redisHealthy {
			if time.Since(lastErrorTime) > 30*time.Second {
				redisHealthy = true // Try to "Half-Open" (simplified here)
			} else {
				next(w, r)
				return
			}
		}

		// Execute the Lua Script atomically(all at once and one execution at a time)
		// KEYS[1] = limitKey
		// ARGV = now, windowStartTime, limit, duration(in seconds)
		result, err := slidingWindowLua.Run(ctx, rdb, []string{limitKey},
			now, windowStartTime, limit, int(duration.Seconds())).Int()
		if err != nil {
			// 1. LOG THE ERROR: So you know Redis is down (Check your terminal)
			fmt.Printf("REDIS ERROR (Failing Open): %v\n", err)
			redisHealthy = false
			lastErrorTime = time.Now()
			fmt.Printf("Circuit Breaker Tripped: %v\n", err)
			// 2. FAIL OPEN: Call the next handler and exit the middleware
			// This lets the user see the page even if Redis is dead
			next(w, r)
			return
		}

		if result == 1 {
			http.Error(w, "Atomic Sliding Window Limit Reached!", 429)
			return
		}
		next(w, r)
	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HOME page!.")
}
func loginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the LOGIN page!.")
}
func main() {
	// 1. Establishing connection to local redis server using rdb client.
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	http.HandleFunc("/", rateLimitMiddleware(homePage, 10, 30*time.Second))
	http.HandleFunc("/login", rateLimitMiddleware(loginPage, 3, 1*time.Minute))

	fmt.Println("Server starting at :8080...")
	http.ListenAndServe(":8080", nil) // spin up the server.

}
