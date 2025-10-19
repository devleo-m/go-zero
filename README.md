# 🚀 GO ZERO

<div align="center">

**The Ultimate Full-Stack Learning Playground**

*A deliberately over-engineered project designed to teach you EVERYTHING about modern backend development with Go.*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=for-the-badge)](https://github.com/yourusername/go-zero)

</div>

---

## 🎯 What is GO ZERO?

**GO ZERO** is not a real product. It's a comprehensive learning laboratory where you implement every major backend pattern used in production systems today.

Think of it as your personal playground to master:

| 🛠️ **Core Features** | 📚 **Learning Areas** |
|---------------------|----------------------|
| 💰 Payment processing (Stripe) | Real-world integrations |
| 🎥 Video streaming (HLS) | Advanced Go patterns |
| 💬 Real-time chat (WebSocket) | Production architecture |
| 🛒 E-commerce flows | Performance optimization |
| 🎓 Course platforms | Security best practices |
| 🎫 Ticketing systems | Monitoring & observability |
| 🏆 Gamification | Testing strategies |
| 📊 Analytics & monitoring | DevOps practices |
| 🔐 Advanced authentication | Clean code principles |

All in one codebase, following industry best practices and production-ready patterns.

---

## 🤔 Why This Project Exists?

As a developer learning Go, you face a common problem:

- ❌ **Tutorials are too simple** (just CRUD operations)
- ❌ **Real projects are overwhelming** (thousands of files)
- ❌ **You don't know what "production-ready" means**

**GO ZERO** solves this by providing a structured path to implement:

- ✅ Everything a senior backend engineer should know
- ✅ Real-world integrations (Stripe, S3, WebSockets)
- ✅ Proper architecture (Hexagonal + DDD)
- ✅ Production patterns (caching, queues, observability)

---

## 🏗️ Architecture

We follow **Hexagonal Architecture** (Ports & Adapters) to keep the business logic isolated and testable.

```
┌─────────────────────────────────────────────────────────┐
│                    DELIVERY LAYER                        │
│  (HTTP, WebSocket, CLI, gRPC - External World)          │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│                   APPLICATION LAYER                      │
│         (Use Cases - Orchestration Logic)                │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│                    DOMAIN LAYER                          │
│     (Entities, Value Objects - Business Rules)           │
└────────────────┬────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────┐
│                 INFRASTRUCTURE LAYER                     │
│  (Database, Cache, Queue, Storage - External Services)   │
└─────────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
go-zero/
├── cmd/
│   └── api/
│       └── main.go                    # 🚀 Application entry point
│
├── internal/
│   ├── modules/                       # 📦 Feature modules (bounded contexts)
│   │   ├── ecommerce/
│   │   │   ├── domain/               # Business entities & rules
│   │   │   ├── usecases/             # Application logic
│   │   │   ├── delivery/             # HTTP handlers
│   │   │   └── repository/           # Data access
│   │   │
│   │   ├── courses/                  # Course platform module
│   │   ├── chat/                     # Real-time chat module
│   │   ├── tickets/                  # Support tickets module
│   │   ├── payments/                 # Payment processing module
│   │   ├── gamification/             # Points & achievements module
│   │   └── analytics/                # Analytics & monitoring module
│   │
│   ├── shared/                        # 🔧 Shared utilities
│   │   ├── auth/                     # JWT, OAuth, 2FA
│   │   ├── storage/                  # S3/MinIO client
│   │   ├── cache/                    # Redis wrapper
│   │   ├── queue/                    # Background jobs
│   │   ├── email/                    # Email service
│   │   └── websocket/                # WebSocket hub
│   │
│   └── infrastructure/                # ⚙️ Technical concerns
│       ├── http/                     # Gin server setup
│       ├── persistence/              # Database connections
│       ├── config/                   # Configuration
│       ├── logger/                   # Structured logging (Zap)
│       └── monitoring/               # Prometheus metrics
│
├── migrations/                        # 🗄️ Database migrations
├── scripts/                           # 🛠️ Helper scripts
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## 🎓 Learning Modules

Each module teaches you a critical real-world skill:

### 1️⃣ E-Commerce 🛒

**What you'll learn:**
- Product CRUD with image uploads
- Shopping cart (Redis-based)
- Checkout flow with Stripe
- Inventory management
- Discount coupons
- Order processing

**Tech Stack:** `Gin` • `GORM` • `Redis` • `Stripe SDK` • `MinIO`

---

### 2️⃣ Course Platform 🎓

**What you'll learn:**
- Video upload (multipart/chunked)
- HLS transcoding (ffmpeg)
- Signed URLs for security
- Progress tracking
- Certificate generation (PDF)

**Tech Stack:** `MinIO (S3)` • `ffmpeg` • `gofpdf`

---

### 3️⃣ Real-Time Chat 💬

**What you'll learn:**
- WebSocket server
- Pub/Sub with Redis
- Message persistence (MongoDB)
- Typing indicators
- Online presence
- Notification system

**Tech Stack:** `Gorilla WebSocket` • `Redis Pub/Sub` • `MongoDB`

---

### 4️⃣ Payment Processing 💰

**What you'll learn:**
- Stripe integration
- Webhook handling
- Payment intents
- Refunds
- Subscription billing
- Split payments (marketplace)

**Tech Stack:** `Stripe SDK` • `Webhook verification`

---

### 5️⃣ Ticketing System 🎫

**What you'll learn:**
- Support ticket lifecycle
- File attachments
- Agent assignment
- SLA tracking
- Status management

**Tech Stack:** `GORM` • `MinIO` • `PostgreSQL`

---

### 6️⃣ Gamification 🏆

**What you'll learn:**
- Points system
- Badges/Achievements
- Leaderboards
- Daily missions
- Reward distribution

**Tech Stack:** `PostgreSQL` • `Redis (sorted sets)`

---

### 7️⃣ Authentication 🔐

**What you'll learn:**
- JWT (access + refresh tokens)
- OAuth 2.0 (Google, GitHub)
- Two-Factor Authentication (TOTP)
- Magic links (passwordless)
- Role-Based Access Control (RBAC)

**Tech Stack:** `jwt-go` • `oauth2` • `otp`

---

### 8️⃣ Background Jobs ⚙️

**What you'll learn:**
- Async task processing
- Job queues (Redis-based)
- Retry mechanisms
- Cron jobs
- Worker pools

**Tech Stack:** `Asynq` • `Redis`

---

### 9️⃣ Performance ⚡

**What you'll learn:**
- Multi-layer caching strategy
- Rate limiting (per user/IP)
- Circuit breaker pattern
- Connection pooling
- Query optimization

**Tech Stack:** `Redis` • `golang-lru` • `Circuit breaker`

---

### 🔟 Observability 📊

**What you'll learn:**
- Structured logging (Zap)
- Metrics (Prometheus)
- Distributed tracing (Jaeger)
- Health checks
- Alerting

**Tech Stack:** `Zap` • `Prometheus` • `Grafana` • `Jaeger`

---

## 🛠️ Tech Stack

| **Category** | **Technology** |
|-------------|----------------|
| **Language** | Go 1.21+ |
| **HTTP Framework** | Gin |
| **WebSocket** | Gorilla WebSocket |
| **Database** | PostgreSQL (primary), MongoDB (chat) |
| **Cache/Queue** | Redis |
| **Object Storage** | MinIO (S3-compatible) |
| **ORM** | GORM |
| **Payments** | Stripe |
| **Video Processing** | ffmpeg |
| **Auth** | JWT, OAuth 2.0, TOTP |
| **Email** | SendGrid / Mailhog (dev) |
| **Logging** | Zap |
| **Metrics** | Prometheus + Grafana |
| **Testing** | Testify, Mockery |
| **Migration** | golang-migrate |
| **Background Jobs** | Asynq |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-zero.git
cd go-zero

# Copy environment variables
cp .env.example .env

# Start infrastructure services
docker-compose up -d

# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start the server
make run
```

The API will be available at `http://localhost:8080`

### Available Services

After running `docker-compose up -d`:

| **Service** | **URL** | **Description** |
|-------------|---------|-----------------|
| **PostgreSQL** | `localhost:5432` | Primary database |
| **Redis** | `localhost:6379` | Cache & queue |
| **MongoDB** | `localhost:27017` | Chat messages |
| **MinIO** | `localhost:9000` | Object storage |
| **MinIO Console** | `localhost:9001` | Storage management |
| **Mailhog** | `localhost:8025` | Email testing |
| **Prometheus** | `localhost:9090` | Metrics |
| **Grafana** | `localhost:3000` | Dashboards |

---

## 📚 API Documentation

### Authentication Endpoints

| **Method** | **Endpoint** | **Description** |
|------------|--------------|-----------------|
| `POST` | `/api/v1/auth/register` | Register new user |
| `POST` | `/api/v1/auth/login` | Login (returns JWT) |
| `POST` | `/api/v1/auth/refresh` | Refresh access token |
| `POST` | `/api/v1/auth/logout` | Logout |
| `POST` | `/api/v1/auth/forgot-password` | Request password reset |
| `POST` | `/api/v1/auth/oauth/google` | OAuth login |
| `POST` | `/api/v1/auth/2fa/enable` | Enable 2FA |
| `POST` | `/api/v1/auth/2fa/verify` | Verify 2FA code |

### E-Commerce Endpoints

| **Method** | **Endpoint** | **Description** |
|------------|--------------|-----------------|
| `GET` | `/api/v1/products` | List products |
| `POST` | `/api/v1/products` | Create product (admin) |
| `GET` | `/api/v1/products/:id` | Get product details |
| `PUT` | `/api/v1/products/:id` | Update product (admin) |
| `DELETE` | `/api/v1/products/:id` | Delete product (admin) |
| `POST` | `/api/v1/cart/add` | Add to cart |
| `GET` | `/api/v1/cart` | Get cart |
| `DELETE` | `/api/v1/cart/:item_id` | Remove from cart |
| `POST` | `/api/v1/checkout` | Create checkout session |
| `POST` | `/api/v1/webhooks/stripe` | Stripe webhook handler |

### Course Endpoints

| **Method** | **Endpoint** | **Description** |
|------------|--------------|-----------------|
| `GET` | `/api/v1/courses` | List courses |
| `POST` | `/api/v1/courses` | Create course (instructor) |
| `GET` | `/api/v1/courses/:id` | Get course details |
| `POST` | `/api/v1/courses/:id/enroll` | Enroll in course |
| `POST` | `/api/v1/lessons/:id/upload` | Upload video |
| `GET` | `/api/v1/lessons/:id/stream` | Get streaming URL |
| `POST` | `/api/v1/lessons/:id/progress` | Update progress |

### Chat Endpoints

| **Method** | **Endpoint** | **Description** |
|------------|--------------|-----------------|
| `WS` | `/ws/chat` | WebSocket connection |
| `POST` | `/api/v1/messages` | Send message |
| `GET` | `/api/v1/messages/:room_id` | Get message history |

### Ticket Endpoints

| **Method** | **Endpoint** | **Description** |
|------------|--------------|-----------------|
| `POST` | `/api/v1/tickets` | Create ticket |
| `GET` | `/api/v1/tickets` | List tickets |
| `GET` | `/api/v1/tickets/:id` | Get ticket details |
| `PUT` | `/api/v1/tickets/:id` | Update ticket |
| `POST` | `/api/v1/tickets/:id/messages` | Add message |

> 📖 **Full API documentation** available at `/swagger` (when running)

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run e2e tests
make test-e2e

# Generate mocks
make mock
```

### Test Structure

```
internal/
└── modules/
    └── ecommerce/
        ├── domain/
        │   └── product_test.go        # Unit tests (business logic)
        ├── usecases/
        │   └── create_product_test.go # Use case tests (with mocks)
        └── delivery/
            └── product_handler_test.go # Integration tests
```

---

## 📊 Performance Benchmarks

Run benchmarks with:
```bash
make bench
```

### Target Metrics (on local machine)

| **Metric** | **Target** |
|------------|------------|
| **Latency p50** | < 10ms |
| **Latency p99** | < 50ms |
| **Throughput** | > 10k req/s (simple endpoints) |
| **Cache hit ratio** | > 80% |
| **Database connection pool** | 95%+ utilization |

---

## 🛡️ Security Features

- ✅ **JWT with refresh token rotation**
- ✅ **Password hashing (bcrypt)**
- ✅ **Rate limiting (per IP and per user)**
- ✅ **CORS configuration**
- ✅ **SQL injection prevention** (GORM parameterized queries)
- ✅ **XSS protection** (sanitized inputs)
- ✅ **CSRF tokens** (for forms)
- ✅ **Helmet-like headers** (security headers)
- ✅ **Input validation** (go-playground/validator)
- ✅ **File upload restrictions** (size, type)
- ✅ **Signed URLs for private content**
- ✅ **2FA support**

---

## 📈 Monitoring & Observability

### Metrics (Prometheus)
Access metrics at `http://localhost:8080/metrics`

**Custom metrics:**
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `cache_hits_total` - Cache hit/miss ratio
- `background_jobs_processed` - Job processing stats
- `websocket_connections` - Active WebSocket connections

### Logs (Zap)
Structured JSON logs with:
- Request ID (for tracing)
- User ID
- Timestamp
- Level (debug, info, warn, error)
- Context fields

### Tracing (Jaeger)
Distributed tracing for:
- HTTP requests
- Database queries
- Cache operations
- External API calls

---

## 🗺️ Roadmap

### Phase 1: Foundation ✅
- [x] Hexagonal architecture setup
- [x] Docker environment
- [x] Database migrations
- [x] Configuration management
- [x] Logging & monitoring

### Phase 2: Core Features (In Progress)
- [ ] Authentication (JWT)
- [ ] User management
- [ ] E-commerce module
- [ ] Payment integration
- [ ] File uploads

### Phase 3: Advanced Features
- [ ] Course platform
- [ ] Video streaming
- [ ] Real-time chat
- [ ] Background jobs
- [ ] Email notifications

### Phase 4: Performance
- [ ] Multi-layer caching
- [ ] Rate limiting
- [ ] Circuit breaker
- [ ] Load testing

### Phase 5: Production Ready
- [ ] 80%+ test coverage
- [ ] CI/CD pipeline
- [ ] Docker production image
- [ ] Kubernetes manifests
- [ ] Comprehensive docs

---

## 🤝 Contributing

This is a learning project, but contributions are welcome!

1. **Fork** the repository
2. **Create** your feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add some amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Please ensure:
- ✅ Tests pass (`make test`)
- ✅ Code follows Go conventions (`make lint`)
- ✅ Documentation is updated
- ✅ Commit messages are clear

---

## 📝 Learning Resources

Want to understand the concepts better?

- **Hexagonal Architecture:** [Alistair Cockburn's original article](https://alistair.cockburn.us/hexagonal-architecture/)
- **Domain-Driven Design:** [DDD in Go](https://github.com/marcusolsson/gophercon-eu-2018)
- **Go Best Practices:** [Effective Go](https://golang.org/doc/effective_go.html)
- **Stripe Integration:** [Stripe Go SDK](https://stripe.com/docs/api/go)
- **WebSockets in Go:** [Gorilla WebSocket](https://github.com/gorilla/websocket)

---

## 🎯 Project Philosophy

> **"Learn by doing, not by reading."**

This project follows these principles:

- **Progressive Complexity** - Start simple, add features incrementally
- **Production Patterns** - No toy code, everything is "real world"
- **Clean Architecture** - Business logic independent of frameworks
- **Testable** - Every layer can be tested in isolation
- **Observable** - You can see what's happening inside
- **Documented** - Code explains itself + comprehensive docs

---

## 📜 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by real-world production systems
- Built to help developers level up their Go skills
- Community-driven learning approach

---

## 💬 Contact

Have questions? Found a bug? Want to contribute?

- **GitHub Issues:** [Create an issue](https://github.com/yourusername/go-zero/issues)
- **Discussions:** [Join the conversation](https://github.com/yourusername/go-zero/discussions)
- **Email:** your.email@example.com

---

<div align="center">

**Made with ❤️ for developers learning Go**

⭐ **Star this repo if you find it helpful!**

[Report Bug](https://github.com/yourusername/go-zero/issues) · [Request Feature](https://github.com/yourusername/go-zero/issues) · [Contribute](https://github.com/yourusername/go-zero/pulls)

</div>
