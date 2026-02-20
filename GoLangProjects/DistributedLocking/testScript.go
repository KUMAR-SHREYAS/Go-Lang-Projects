package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func RunChaosTest() {
	// Setup Redis Client
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	lockKey := "chaos_test_lock"

	// Cleanup any old locks from previous runs
	rdb.Del(lockKey)

	fmt.Println("🧪 Starting Chaos Test for Distributed Locker...")

	// --- TEST 1: Expiration Safety ---
	fmt.Println("\n--- Test 1: Expiration Safety ---")
	token1 := uuid.New().String()
	success1 := AcquireLock(rdb, lockKey, token1, 2*time.Second)
	
	if success1 {
		fmt.Println("Worker 1 acquired lock for 2s...")
	}

	fmt.Println("Worker 1 is sleeping past the expiration...")
	time.Sleep(3 * time.Second) // Wait for lock to expire in Redis

	// Worker 2 tries to grab the now-expired lock
	token2 := uuid.New().String()
	success2 := AcquireLock(rdb, lockKey, token2, 5*time.Second)
	if success2 {
		fmt.Println("✅ Success: Worker 2 grabbed the expired lock.")
	}

	// --- TEST 2: Token Mismatch Protection (The "Stealing" Test) ---
	fmt.Println("\n--- Test 2: Token Mismatch Protection ---")
	fmt.Println("Worker 1 wakes up and tries to release the lock it no longer owns...")
	
	// Worker 1 still thinks its token1 is valid, but the lock is now held by token2
	released := ReleaseLock(rdb, lockKey, token1) 
	
	if !released {
		fmt.Println("✅ Success: Lua script blocked Worker 1 from deleting Worker 2's lock!")
	} else {
		fmt.Println("❌ Critical Failure: Worker 1 successfully deleted Worker 2's lock!")
	}

	// --- TEST 3: Correct Release ---
	fmt.Println("\n--- Test 3: Valid Release ---")
	released2 := ReleaseLock(rdb, lockKey, token2)
	if released2 {
		fmt.Println("✅ Success: Worker 2 released its own lock correctly.")
	}

	// Final cleanup
	rdb.Del(lockKey)
	fmt.Println("\n🧪 Chaos Test Suite Complete.")
}