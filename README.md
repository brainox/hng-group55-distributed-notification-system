<div align="center">

# 🔔 Distributed Notification System

### HNG Group 55 - Scalable Microservices Architecture

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Microservices](https://img.shields.io/badge/architecture-microservices-orange)
![Docker](https://img.shields.io/badge/docker-enabled-2496ED?logo=docker)

</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Core Services](#-core-services)
- [Infrastructure](#-infrastructure)
- [Observability](#-observability)
- [Getting Started](#-getting-started)
- [Documentation](#-documentation)

---

## 🎯 Overview

A highly scalable, distributed notification system built with microservices architecture. This system handles multi-channel notifications (email, push, SMS) with high throughput, reliability, and observability.

### ✨ Key Features

- 🚀 **High Performance** - Handles thousands of notifications per second
- 🔄 **Multi-Channel Support** - Email, Push notifications, and more
- 🛡️ **Resilience** - Built-in circuit breakers, retries, and idempotency
- 📊 **Full Observability** - Metrics, logs, and distributed tracing
- 🔌 **Event-Driven** - Kafka-based message streaming
- 🐳 **Container-Ready** - Fully dockerized for easy deployment

---

## 🏗️ Architecture

This system follows a microservices architecture pattern with event-driven communication, leveraging Apache Kafka for reliable message streaming and Redis for caching and session management.

---

## 📁 Project Structure

```
notification-system/
│
├── 🌐 api_gateway/              # API Gateway & routing layer
│
├── ⚙️  services/
│   ├── user_service/            # User management & authentication
│   ├── template_service/        # Notification template management
│   ├── email_service/           # Email notification handler
│   ├── push_service/            # Push notification handler
│
├── 🏢 infra/
│   ├── kafka/                   # Message broker configuration
│   ├── redis/                   # Caching & session store
│   ├── postgres/                # Primary database
│   ├── nginx/                   # Load balancer & reverse proxy
│
├── 🔧 shared/
│   └── libs/
│       ├── circuit_breaker/     # Circuit breaker pattern
│       ├── idempotency/         # Idempotency handling
│       ├── retry/               # Retry logic & backoff
│       └── logging/             # Centralized logging utilities
│
├── 📊 observability/
│   ├── prometheus/              # Metrics collection
│   ├── grafana/                 # Metrics visualization
│   ├── loki/                    # Log aggregation
│   ├── jaeger/                  # Distributed tracing
│   └── alerting/                # Alert management
│
├── 🚀 deployments/
│   ├── docker/                  # Docker compose configurations
│   ├── staging/                 # Staging environment configs
│   └── production/              # Production environment configs
│
├── 🔄 .github/
│   └── workflows/               # CI/CD pipelines
│
└── 📚 docs/
    ├── architecture_diagram/    # System architecture diagrams
    ├── openapi_specs/           # API specifications
    └── readmes/                 # Additional documentation
```

---

## 🔌 Core Services

| Service | Description | Port |
|---------|-------------|------|
| **API Gateway** | Entry point for all client requests, handles routing and authentication | `8000` |
| **User Service** | Manages user accounts, preferences, and authentication | `8001` |
| **Template Service** | Handles notification templates and personalization | `8002` |
| **Email Service** | Processes and sends email notifications | `8003` |
| **Push Service** | Handles push notifications to mobile devices | `8004` |

---

## 🏢 Infrastructure

### Message Broker
- **Apache Kafka** - Event streaming platform for real-time data pipelines

### Data Storage
- **PostgreSQL** - Primary relational database
- **Redis** - High-performance caching and session management

### Load Balancing
- **Nginx** - Reverse proxy and load balancer

---

## 📊 Observability

### Monitoring Stack

| Tool | Purpose |
|------|---------|
| **Prometheus** | Metrics collection and alerting |
| **Grafana** | Metrics visualization and dashboards |
| **Loki** | Log aggregation and querying |
| **Jaeger** | Distributed request tracing |

### Key Metrics Tracked
- Request throughput and latency
- Service health and uptime
- Queue depths and processing rates
- Error rates and types
- Resource utilization (CPU, memory, disk)

---

## 🚀 Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js (v18+) or Python (v3.10+)
- Kafka & Zookeeper
- PostgreSQL
- Redis

### Quick Start

```bash
# Clone the repository
git clone https://github.com/brainox/hng-group55-distributed-notification-system.git
cd hng-group55-distributed-notification-system

# Start infrastructure services
docker-compose -f deployments/docker/docker-compose.yml up -d

# Start individual services (example)
cd services/email_service
npm install && npm start
```

### Environment Setup

Each service requires its own environment configuration. Copy the example env files:

```bash
cp .env.example .env
```

Update the `.env` files with your specific configuration values.

---

## 📚 Documentation

- [Architecture Overview](docs/architecture_diagram/)
- [API Documentation](docs/openapi_specs/)
- [Service READMEs](docs/readmes/)
- [Deployment Guide](deployments/)

---

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

---

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 👥 Team

**HNG Group 55** - Building scalable notification systems

---

<div align="center">

Made with ❤️ by HNG Group 55

</div>