# 🛡️ Distributed Locking (Redis-Backed)

A production-grade implementation of a **Distributed Mutual Exclusion (Mutex)** lock in Go. This pattern ensures that in a cluster of multiple servers, only one instance can perform a specific "critical" operation at any given time.

---

## 🛠️ Key Features

* **Mutual Exclusion (`SETNX`)**: Uses Redis `SET if Not Exists` logic to ensure a lock is granted only to the first requester.
* **Owner Identification (UUID Tokens)**: Every lock attempt generates a unique **UUID**. This serial number acts as a "security badge," ensuring that one server cannot release a lock held by another.
* **Atomic Safe-Release (Lua)**: To prevent race conditions during the unlock process, we use a **Lua script**. It verifies the `Owner Token` matches the value in Redis before performing the `DEL` operation.
* **Deadlock Prevention**: Integrated **TTL (Time-To-Live)** ensures that if a server crashes while holding a lock, the lock is automatically released after a safety timeout.
* **Acquisition Retry Loop**: Includes a configurable retry mechanism with **Exponential Backoff**, allowing servers to "wait in line" rather than failing immediately.

---

## 🚦 How It Works



1.  **Request**: Server A tries to set a key in Redis with a unique UUID and a 10s expiration.
2.  **Logic Execution**: Server A performs the critical task (e.g., generating an invoice).
3.  **Safety Check**: If Server A takes 11s, the lock expires. Server B grabs it.
4.  **Atomic Release**: When Server A tries to unlock, the Lua script sees that the token in Redis belongs to Server B and **refuses to delete it**.

---

## 💻 Technical Competencies Demonstrated

* **Go Concurrency**: Managing context and timing for retries and backoffs.
* **Redis Internals**: Using `SetNX` for atomicity and Lua scripts for complex server-side logic.
* **System Reliability**: Implementing "Owner ID" checks to prevent cascading failures in distributed environments.



---

## 🚀 Quick Start

1. **Install Dependency**:
   ```bash
   go get [github.com/google/uuid](https://github.com/google/uuid)