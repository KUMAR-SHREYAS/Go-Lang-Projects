package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var ctx = context.Background()

// AcquireLock tries to grab a unique lock in Redis
func AcquireLock(rdb *redis.Client, lockKey string, token string, expiration time.Duration) bool {

	// 1. SET if Not Exists (NX) with an Expiration (EX)
	success, err := rdb.SetNX(lockKey, token, expiration).Result()
	if err != nil {
		log.Printf("Redis Error: %v", err)
		return false
	}

	return success
}

// Lua script: Returns 1 if deleted, 0 if token mismatch or key doesn't exist
var unlockLua = redis.NewScript(`
	if redis.call("get", KEYS[1])==ARGV[1] then
		return redis.call("del", KEYS[1])
	else 
		return 0
	end
`)

// ReleaseLock uses Lua to ensure we only delete our OWN lock
func ReleaseLock(rdb *redis.Client, lockKey string, token string) bool {
	result, err := unlockLua.Run(rdb, []string{lockKey}, token).Int()
	if err != nil {
		log.Printf("unlock Error: %v", err)
		return false
	}
	return result == 1
}
func AcquireRedlock(nodes []*redis.Client, key string, expiration time.Duration) (string, bool) {
	n := len(nodes)
	quorum := (n / 2) + 1
	token := uuid.New().String()

	successCount := 0
	startTime := time.Now()

	//1. Try to lock all nodes
	for _, node := range nodes {
		// we use a very short timeout for each individual network call
		if AcquireLock(node, key, token, expiration) {
			successCount++
		}
	}
	//2. Calculate time spent and check Quorum
	elapsed := time.Since(startTime)
	//Subtract elapsed time from TTL to get a actual "valid" time remaining
	validityTime := expiration-elapsed

	if successCount >= quorum && validityTime >0 {
		return token, true
	}

	//3. FAILURE: IF we didn't get quorum , unlock EVERYTHING immediately we don't want locked hanging nodes
	fmt.Println("⚠️ Failed to reach Quorum. Cleaning up...")
	for _, node := range nodes {
		ReleaseLock(node, key, token)
	}
	return "", false
}

func main() {
	// Simulate 3 independent Redis instances (e.g., ports 6379, 6380, 6381)
	// In a real system, these would be different IP addresses.
	node1 := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	node2 := redis.NewClient(&redis.Options{Addr: "localhost:6380"})
	node3 := redis.NewClient(&redis.Options{Addr: "localhost:6381"})

	redisNodes := []*redis.Client{node1, node2, node3}

	lockKey := "global_distributed_lock"

	ttl := 10*time.Second

	token, ok := AcquireRedlock(redisNodes, lockKey, ttl)

	if ok {
		fmt.Printf("✅ Global Lock Acquired! Token: %s\n", token)
		// Perform critical task...
		time.Sleep(2 * time.Second)

		// Release from all nodes
		for _, node := range redisNodes {
			ReleaseLock(node, lockKey, token)
		}
		fmt.Println("🔓 Global Lock Released.")
	}else{
		fmt.Println("❌ Failed to acquire Global Lock.")
	}

	RunChaosTest()
}

