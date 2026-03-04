#  Scalable Go API Gateway

![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![Docker](https://img.shields.io/badge/Docker-Ready-green)
![License](https://img.shields.io/badge/License-MIT-yellow)

A high-performance, cloud-native API Gateway built from scratch in Go. This project acts as a centralized control plane for microservices, handling **Routing**, **Security**, **Rate Limiting**, and **Observability** at the edge.



---

##  Key Features

* **Dynamic Reverse Proxy:** Routes traffic to multiple backend services based on path-prefixes defined in `config.yaml`.
* **IP-Based Rate Limiting:** Implements a thread-safe **Token Bucket algorithm** using `sync.Mutex` to prevent DDoS and "noisy neighbor" issues.
* **Centralized Authentication:** Middleware-based "Security Perimeter" to validate Bearer tokens before requests hit internal services.
* **Distributed Tracing:** Automatic **Request ID (UUID)** injection into headers for end-to-end log correlation.
* **Production Observability:** Built-in **Prometheus** metrics exporter at the `/metrics` endpoint.
* **Resiliency:** Implements **Graceful Shutdown** using OS signal listeners to ensure zero-downtime deployments.

---

##  Technical Architecture

### 1. The Middleware Chain
The gateway uses a "Decorator Pattern" to wrap the core reverse proxy. Every request passes through a strictly ordered pipeline:
1.  **RequestID:** Assigns a unique trace ID for the request lifecycle.
2.  **RateLimiter:** Throttles clients based on IP addresses using a concurrent-safe map.
3.  **Logger:** Captures request metadata and latency.
4.  **Authenticator:** (Optional) Validates `Authorization` headers for protected routes.



### 2. Concurrency Control
To handle high traffic, the gateway leverages Go's lightweight Goroutines. The rate limiter is designed for thread-safety, using a `sync.Mutex` to protect the bucket state during parallel updates, preventing race conditions.

---

##  Quick Start

### Prerequisites
* **Go 1.24+**
* **Docker & Docker Compose** (for containerized deployment)

### 1. Clone & Run Locally
```bash
git clone [https://github.com/sachinvivek31/Scalable-Api-Gateway.git](https://github.com/sachinvivek31/Scalable-Api-Gateway.git)
cd Scalable-Api-Gateway
go run main.go
```

### 2. Run with Docker
```bash
docker compose up --build
```

## Testing
**Feature ** |  **Command**                                                                              | **Expected Result**
Public Route |  curl.exe -i -X POST http://localhost:8080/post                                           | 200 OK (via httpbin)
Secure Route |  curl.exe -i http://localhost:8080/get                                                    | 401 Unauthorized
Auth Success |  "curl.exe -i -H ""Authorization: Bearer my-secret-pro-token"" http://localhost:8080/get" | 200 OK
Metrics      |  curl.exe http://localhost:8080/metrics                                                   | Prometheus Metrics Stream
