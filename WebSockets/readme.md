
# Real-Time Chat Application using WebSockets in Go

A **real-time chat application** built with **Go (Golang) WebSockets** and a **vanilla JavaScript frontend**.
The system demonstrates how to build a **scalable, secure, and event-driven WebSocket server** supporting chatrooms, authentication, and encrypted communication.

This project follows the architecture demonstrated in Percy Bolmér's WebSocket tutorial and extends it into a **complete working system with production-grade concepts** such as:

* Client management
* Event-driven messaging
* Ping/Pong heartbeat
* Message size limits
* Origin security checks
* OTP-based authentication
* Secure communication via **HTTPS/WSS**

---

# Architecture Overview

```
Client (Browser)
      │
      │  WebSocket (ws / wss)
      ▼
Go WebSocket Server
      │
      ├── Manager
      │     ├── Client Registry
      │     ├── Event Handlers
      │     └── OTP Authentication
      │
      ├── Client
      │     ├── readMessages()
      │     └── writeMessages()
      │
      └── Event System
            ├── send_message
            ├── new_message
            └── change_room
```

---

# Features

### Real-Time Messaging

* Full duplex communication using **WebSockets**
* Messages broadcast to clients in the same chatroom

### Chatrooms

Users can dynamically switch chatrooms.

```
general
sports
movies
tech
```

Messages are only broadcast to clients within the same room.

---

### Event Driven WebSocket API

Messages follow a structured **JSON event format**.

Example:

```json
{
  "type": "send_message",
  "payload": {
    "message": "Hello",
    "from": "user"
  }
}
```

Supported events:

| Event          | Description              |
| -------------- | ------------------------ |
| `send_message` | Send a chat message      |
| `new_message`  | Broadcasted chat message |
| `change_room`  | Switch chatrooms         |

---

### Client Manager

The `Manager` maintains all connected clients.

Responsibilities:

* Accept WebSocket connections
* Register / remove clients
* Route incoming events
* Broadcast messages

---

### Concurrent Message Handling

Each client runs two goroutines:

```
Client
 ├── readMessages()
 └── writeMessages()
```

An **egress channel** prevents concurrent WebSocket writes.

```
egress chan Event
```

This follows Gorilla WebSocket best practices.

---

### Heartbeat (Ping / Pong)

To keep connections alive:

Server sends periodic **Ping** frames.

Client automatically responds with **Pong**.

```
Ping Interval: 9 seconds
Pong Timeout: 10 seconds
```

If a client stops responding, the server disconnects it.

---

### Security Features

#### Message Size Limits

Prevents malicious payloads.

```
Max message size: 512 bytes
```

Configured with:

```
SetReadLimit(512)
```

---

#### Origin Checking

Prevents **Cross-Site WebSocket Hijacking**.

Example allowed origin:

```
https://localhost:8080
```

---

#### OTP Authentication

Users authenticate through an HTTP endpoint.

Flow:

```
1. User logs in via HTTP
2. Server returns OTP
3. Client connects WebSocket using OTP
```

Example:

```
wss://localhost:8080/ws?otp=TOKEN
```

OTPs expire after **5 seconds**.

---

### HTTPS / WSS Encryption

Communication uses **secure WebSockets (WSS)**.

```
HTTPS + WSS
```

Traffic is encrypted using TLS certificates.

---

# Project Structure

```
WebSockets/
│
├── main.go
├── manager.go
├── client.go
├── event.go
├── otp.go
│
├── frontend/
│     └── index.html
│
├── server.crt
├── server.key
│
├── tmp/
│     └── main.exe
│
└── README.md
```

---

# File Responsibilities

### main.go

Application entry point.

Responsibilities:

* Initialize context
* Start HTTP server
* Register routes

Routes:

```
/        -> Frontend
/ws      -> WebSocket endpoint
/login   -> Authentication
/debug   -> Debug info
```

---

### manager.go

Handles:

* WebSocket upgrades
* Client registry
* Event routing
* OTP validation

---

### client.go

Represents a connected user.

Responsibilities:

```
readMessages()
writeMessages()
pongHandler()
```

Maintains:

```
connection
manager
egress channel
chatroom
```

---

### event.go

Defines the **WebSocket event protocol**.

Structures:

```
Event
SendMessageEvent
NewMessageEvent
ChangeRoomEvent
```

Contains event handlers:

```
SendMessageHandler
ChatRoomHandler
```

---

### otp.go

Implements a **temporary authentication token system**.

Functions:

```
NewOTP()
VerifyOTP()
Retention()
```

Expired tokens are automatically removed.

---

### frontend/index.html

Simple client interface using **vanilla JavaScript**.

Features:

* Login form
* Chat interface
* Chatroom selection
* WebSocket connection
* Event routing

---

# Installation

### 1. Clone Repository

```
git clone https://github.com/yourusername/go-websocket-chat.git

cd go-websocket-chat
```

---

### 2. Install Dependencies

```
go mod tidy
```

Required package:

```
github.com/gorilla/websocket
```

---

### 3. Generate TLS Certificates

Run:

```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365
```

These are **development certificates only**.

---

### 4. Run Server

```
go run *.go
```

Server starts at:

```
https://localhost:8080
```

---

# Usage

### Step 1 — Open Application

Visit:

```
https://localhost:8080
```

Accept the self-signed certificate warning.

---

### Step 2 — Login

Example credentials:

```
Username: percy
Password: 123
```

Server returns an **OTP token**.

---

### Step 3 — Connect WebSocket

Client automatically connects to:

```
wss://localhost:8080/ws?otp=TOKEN
```

---

### Step 4 — Start Chatting

Open two browser tabs to test messaging.

---

# Example WebSocket Events

### Send Message

```
{
  "type": "send_message",
  "payload": {
    "message": "Hello world",
    "from": "user"
  }
}
```

---

### Broadcast Message

```
{
  "type": "new_message",
  "payload": {
    "message": "Hello world",
    "from": "user",
    "sent": "2026-01-01T12:00:00Z"
  }
}
```

---

# Debug Endpoint

Shows number of connected clients.

```
https://localhost:8080/debug
```

---

# Development Notes

Ensure these files are ignored by Git:

```
server.key
server.crt
tmp/
*.exe
```

These contain sensitive or build artifacts.

---

# Future Improvements

Possible enhancements:

* JWT authentication
* Persistent chat history
* Redis pub/sub for scaling
* Kubernetes deployment
* Rate limiting
* User presence indicators
* Message persistence
* Typing indicators

---

# Learning Outcomes

This project demonstrates:

* WebSocket protocol fundamentals
* Concurrent Go networking
* Event-driven architectures
* Real-time application design
* Secure WebSocket communication
* Client-server message routing

---

# Acknowledgements

Based on the excellent tutorial by:

**Percy Bolmér**

Article:

Mastering WebSockets With Go

---

# License

MIT License
