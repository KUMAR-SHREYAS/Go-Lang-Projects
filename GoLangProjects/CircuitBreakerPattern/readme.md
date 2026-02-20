# ⚡ Adaptive Distributed Locker & Circuit Breaker

A high-resilience Go implementation of a **Distributed Mutex** that intelligently distinguishes between network failures and resource contention. This project combines the **Redlock** algorithm with an **Error-Aware Circuit Breaker** to provide a "self-healing" concurrency primitive.

---

## 🏗️ Architecture



### 1. Smart Circuit Breaker (Resilience)
The system categorizes errors into two distinct types to ensure high availability:
* **System Faults (5xx style):** Connection resets or Redis crashes. These **trip the Circuit Breaker** to protect the application from hanging on network timeouts.
* **Logic Contention (4xx style):** "Lock Busy" signals. These are **ignored by the Breaker** because they indicate the system is healthy, just busy. This prevents a popular resource from accidentally shutting down the entire service.

### 2. Distributed Locking (Redlock)
* **Mutual Exclusion:** Uses Redis `SETNX` to ensure only one worker owns a resource at a time.
* **Atomic Release:** Implemented via a **Lua Script**. This ensures a "slow worker" cannot accidentally delete a lock that has already been re-assigned to a different process after a timeout.



### 3. Reliability & Retry Logic
* **Sequential Polling:** Instead of rejecting requests when a lock is busy, workers enter a retry loop to "wait in line."
* **Short-Circuiting:** If the Circuit Breaker is `OPEN`, the retry loop terminates early to avoid wasting CPU cycles on a known-dead connection.

---

## 🚦 Request Flow



1.  **Entry:** Worker requests a lock for a specific key.
2.  **Filter:** Circuit Breaker checks if the path to Redis is healthy.
3.  **Attempt:** If `CLOSED`, the worker attempts a `SETNX` on the Redis node.
4.  **Evaluate:**
    * **Success:** Task executes; Lock is released via Lua.
    * **Lock Busy:** Worker sleeps and retries (Breaker remains `CLOSED`).
    * **Network Error:** Breaker increments failure count and eventually trips to `OPEN`.

---

## 🛠️ Tech Stack & Patterns
* **Language:** Go (Golang)
* **Database:** Redis
* **Concurrency:** `sync.Mutex` (State safety), `sync.WaitGroup` (Goroutine sync), `goroutines` (Parallelism).
* **Patterns:** Redlock, Circuit Breaker, Proxy, Retry/Backoff.

---

## 🚀 Getting Started

### Prerequisites
* Go 1.18+
* Redis server running on `localhost:6379`

### Run the Stress Test
```bash
# Initialize the module
go mod init distributed-locker
go mod tidy

# Run the test
go run main.go