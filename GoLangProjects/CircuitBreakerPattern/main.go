package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// --- CIRCUIT BREAKER (Improved Logic) ---
const (
	stateClosed = iota
	StateHalfOpen
	StateOpen
)

type CircuitBreaker struct {
	mutex            sync.Mutex
	state            int
	failureCount     int
	failureThreshold int
	retryTimeout     time.Duration
	lastFailureTime  time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{state: stateClosed, failureThreshold: threshold, retryTimeout: timeout}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()
	if cb.state == StateOpen && time.Since(cb.lastFailureTime) > cb.retryTimeout {
		cb.state = StateHalfOpen
	}
	if cb.state == StateOpen {
		cb.mutex.Unlock()
		return errors.New("circuit open")
	}
	cb.mutex.Unlock()

	err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		// IMPORTANT: Only count actual Redis/Network errors as failures.
		// If the error is just "lock busy", we don't trip the breaker.
		if err.Error() != "lock busy" {
			cb.failureCount++
			cb.lastFailureTime = time.Now()
			if cb.failureCount >= cb.failureThreshold || cb.state == StateHalfOpen {
				cb.state = StateOpen
				fmt.Println("\n🔴 CIRCUIT TRIPPED: Redis is having connection issues.")
			}
		}
		return err
	}

	cb.state = stateClosed
	cb.failureCount = 0
	return nil
}

// --- REDLOCK WITH RETRY ---

var unlockLua = redis.NewScript(`
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`)

func ProtectedAcquireRedlockWithRetry(cb *CircuitBreaker, nodes []*redis.Client, key string, expiration time.Duration, maxRetries int, workerId int) (string, bool) {
	token := uuid.New().String()
	quorum := (len(nodes) / 2) + 1

	for i := 0; i < maxRetries; i++ {
		successCount := 0
		err := cb.Execute(func() error {
			votes := 0
			for _, node := range nodes {
				ok, err := node.SetNX(key, token, expiration).Result()
				if err == nil && ok {
					votes++
				}
			}
			if votes < quorum {
				// We return a specific string so the Breaker knows NOT to trip
				return errors.New("lock busy") 
			}
			successCount = votes
			return nil
		})

		if err == nil && successCount >= quorum {
			return token, true
		}

		// If the breaker is actually OPEN, stop trying
		if err != nil && err.Error() == "circuit open" {
			return "", false
		}

		fmt.Printf("⏳ [Worker %d] Lock busy, retrying attempt %d...\n",workerId, i+1)
		time.Sleep(400 * time.Millisecond) // Short sleep before retry
	}
	return "", false
}



func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	nodes := []*redis.Client{rdb}
	cb := NewCircuitBreaker(3, 5*time.Second)
	var wg sync.WaitGroup

	fmt.Println("🚀 Starting Fixed Stress Test...")

	for i := 1; i <= 3; i++ { // Start with 3 workers to see the sequence clearly
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Increased retries to 15 to give everyone a chance
			token, ok := ProtectedAcquireRedlockWithRetry(cb, nodes, "transfer_lock", 5*time.Second, 15, i)
			if ok {
				fmt.Printf("[Worker %d] ✅ GOT LOCK!\n", id)
				time.Sleep(1 * time.Second) // Simulate work
				unlockLua.Run(rdb, []string{"transfer_lock"}, token)
				fmt.Printf("[Worker %d] 🔓 Released.\n", id)
			} else {
				fmt.Printf("[Worker %d] ❌ FAILED after all retries.\n", id)
			}
		}(i)
	}
	wg.Wait()
}