# 🖥️ Uptime - Distributed Website Monitoring Service

**Uptime** is a modular, event-driven website uptime monitoring service. It leverages Go's concurrency, Kafka for pub/sub event flow, Redis for fast cache access, and PostgreSQL for durable persistence. It is optimized for scale, performance, and modularity.

---

## 🚀 Features

- 🕒 **Scheduled Monitoring** – Periodic HTTP pings to user-defined URLs.
- 📬 **Event-Driven Architecture** – Kafka is used for distributing monitor results and scheduling tasks.
- 🧵 **Worker Pool** – Efficient concurrent message processing to handle thousands of events per minute.
- 🧠 **PostgreSQL Logging** – Durable persistence of monitor status over time.
- ⚡ **Redis Caching** – Real-time monitor status updates for quick reads.
- 🌐 **REST API** – Expose endpoints for CRUD operations on monitors.
- 🔐 **Environment-based Config** – Easy switching between local and cloud environments.

---

## 🛠️ Tech Stack

- **Language**: Go 1.24.3+
- **Message Queue**: Kafka (Confluent v2 client)
- **Database**: PostgreSQL
- **Cache**: Redis
- **API**: Standard Go `net/http`
- **Others**: Docker (for deployment), systemd/graceful shutdown support

---

## 📁 Project Structure

```bash
.
├── cmd/                  # Application entry point
├── internal/
│   ├── cache/            # Redis client and utilities
│   ├── config/           # Environment variable loading
│   ├── constants/        # Shared constant values
│   ├── events/           # Kafka consumer/producer wrapper
│   ├── http/             # HTTP server and handlers
│   ├── logger/           # Kafka consumer for log processing (with worker pool)
│   ├── models/           # Data models
│   ├── postgres/         # Postgres logic
│   ├── scheduler/        # Monitor scheduling engine
│   └── worker/           # Monitor runner and async job handler
├── .env.example          # Environment config template
├── go.mod / go.sum       # Go dependencies
└── README.md             # Project documentation
```
---

**Built with ❤️ by Siimas.**
