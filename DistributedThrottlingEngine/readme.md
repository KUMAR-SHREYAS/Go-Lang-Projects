# 🛡️ Atomic Distributed Rate Limiter (Go + Redis)

A high-performance, resilient middleware designed to protect APIs from abuse while guaranteeing **100% uptime** through failure-tolerance patterns.

---

## 🚀 Key Features

* **Sliding Window Algorithm**: Uses Redis Sorted Sets (**ZSET**) and microsecond precision to provide perfect accuracy, preventing "window-edge bursting" common in fixed-window systems.
* **Atomic Lua Scripting**: Logic is offloaded to the Redis server via **Lua scripts**, reducing network round-trips from 3 to 1 and eliminating **Race Conditions** in distributed environments.
* **Granular Partitioning**: Implements independent "buckets" per **User-IP and Route-Path**, allowing for different security levels (e.g., strict for `/login`, relaxed for `/home`).
* **Self-Healing Resilience**:
    * **Fail-Open**: Ensures the application stays online even if Redis crashes.
    * **Circuit Breaker**: Automatically detects database outages and "trips" to prevent latency spikes, allowing the system to self-heal (Half-Open state) once Redis recovers.

---

## 🛠️ Technical Competencies Demonstrated

* **Golang**: Middleware design, context management, and high-concurrency patterns.
* **Redis**: Advanced data structures (ZSETs), TTL management, and Lua integration.
* **System Design**: Resilience patterns (Circuit Breaker), Atomicity, and Distributed State Management.

---

## 📈 Architecture Detail



1. **Request Ingress**: Middleware captures IP and Path.
2. **Atomic Execution**: A Lua script executes `ZREMRANGEBYSCORE`, `ZCARD`, and `ZADD` in a single Redis heartbeat.
3. **Resilience Layer**: The Circuit Breaker monitors Redis health; if a failure is detected, it enters an **Open** state to protect system latency, eventually transitioning through **Half-Open** to self-heal.