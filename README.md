# Go Microservices - Rental Property Platform

A distributed microservices-based rental property listing application built with Go. This platform enables users to browse, list, and manage rental properties with features like user authentication, favorites, and email notifications.

<table>
  <tr>
    <td width="55%" align="center">
      <img src="./public/home_page.png" alt="home page">
    </td>
	<td align="center">
      <img src="./public/verify_modal.png" alt="verify modal">
      <img src="./public/post.png" alt="post">
    </td>

  </tr>
</table>

## Architecture Overview

```
                                    ┌─────────────────┐
                                    │    Frontend     │
                                    │   (Go/HTML)     │
                                    └────────┬────────┘
                                             │
                                             ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                          Broker Service (API Gateway)                       │
│                                  :8080                                      │
└───────┬──────────┬──────────┬──────────┬──────────┬──────────┬─────────────┘
        │          │          │          │          │          │
        ▼          ▼          ▼          ▼          ▼          ▼
   ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
   │  Auth  │ │  Post  │ │Favorite│ │ Logger │ │  Mail  │ │Listener│
   │Service │ │Service │ │Service │ │Service │ │Service │ │Service │
   │ :8081  │ │        │ │        │ │        │ │        │ │        │
   └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └────────┘ └───┬────┘
       │          │          │          │                      │
       ▼          ▼          ▼          ▼                      ▼
  ┌─────────┐ ┌─────────┐         ┌─────────┐            ┌──────────┐
  │PostgreSQL││PostgreSQL│        │ MongoDB │            │ RabbitMQ │
  │ + Redis │ │         │         │         │            │          │
  └─────────┘ └─────────┘         └─────────┘            └──────────┘
```

## Technology Stack

| Category         | Technology                              |
| ---------------- | --------------------------------------- |
| Language         | Go 1.18 - 1.25                          |
| Web Framework    | Chi v5                                  |
| Databases        | PostgreSQL 14.0, MongoDB 4.2, Redis 8.4 |
| Message Queue    | RabbitMQ 3.9                            |
| Authentication   | JWT, OAuth2 (Google)                    |
| Containerization | Docker, Docker Compose                  |
| Email            | SMTP with go-simple-mail                |

## Services

### Broker Service (API Gateway)

Central entry point that routes all client requests to appropriate backend services.

- **Port**: 8080
- Handles CORS, request forwarding, and OAuth routes

### Authentication Service

Manages user authentication and authorization.

- **Port**: 8081
- JWT token management (access + refresh tokens)
- User registration with email verification
- Password reset flow
- Google OAuth2 integration
- Session management via Redis

### Post Service

Manages rental property listings.

- CRUD operations for posts
- Geolocation-based search (lat/lng with radius)
- Property details (bedrooms, bathrooms, price, images)

### Favourite Service

Manages user favorites/wishlist.

- Add/remove favorites
- Bulk sync from localStorage

### Logger Service

Centralized logging to MongoDB.

- Collects logs from all services
- Log levels: INFO, WARNING, ERROR

### Mail Service

Handles all email communications.

- Email verification
- Password reset emails
- HTML email templates

### Listener Service

Consumes messages from RabbitMQ.

- Processes log events
- Triggers email notifications
- Handles async operations

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.18+ (for local development)
- Make (optional)

### Running with Docker Compose

**Windows:**

```bash
cd project
make up_build
```

**Linux:**

```bash
cd linux_project
make up_build
```

### Manual Docker Compose

```bash
cd project
docker-compose up -d --build
```

### Stopping Services

```bash
make down
```

## API Endpoints

### Authentication

| Method | Endpoint                | Description            |
| ------ | ----------------------- | ---------------------- |
| POST   | `/auth/login`           | User login             |
| POST   | `/auth/register`        | User registration      |
| POST   | `/auth/verify-email`    | Verify email address   |
| POST   | `/auth/forgot-password` | Request password reset |
| POST   | `/auth/reset-password`  | Reset password         |
| GET    | `/auth/profile`         | Get user profile       |
| GET    | `/oauth/google/login`   | Google OAuth login     |

### Posts

| Method | Endpoint                   | Description         |
| ------ | -------------------------- | ------------------- |
| GET    | `/posts`                   | Get all posts       |
| GET    | `/posts/{id}`              | Get single post     |
| GET    | `/posts/author/{authorId}` | Get posts by author |
| POST   | `/posts`                   | Create new post     |
| PUT    | `/posts/{id}`              | Update post         |
| DELETE | `/posts/{id}`              | Delete post         |

### Favorites

| Method | Endpoint                       | Description             |
| ------ | ------------------------------ | ----------------------- |
| GET    | `/favorites/{userId}/ids`      | Get user's favorite IDs |
| POST   | `/favorites`                   | Add to favorites        |
| DELETE | `/favorites/{userId}/{postId}` | Remove from favorites   |
| POST   | `/favorites/sync`              | Bulk sync favorites     |

## Environment Variables

### Authentication Service

```env
DSN=postgres://user:password@postgres:5432/users
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
ACCESS_SECRET=your_access_secret
REFRESH_SECRET=your_refresh_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/oauth/google/callback
```

### Mail Service

```env
MAIL_DOMAIN=your_domain
MAIL_HOST=smtp.your-provider.com
MAIL_PORT=587
MAIL_ENCRYPTION=tls
MAIL_USERNAME=your_username
MAIL_PASSWORD=your_password
FROM_NAME=Your App Name
FROM_ADDRESS=noreply@your-domain.com
```

### Post & Favourite Services

```env
DSN=postgres://user:password@postgres:5432/posts
```

## Project Structure

```
go-micro/
├── broker-service/          # API Gateway
├── authentication-service/  # Auth & user management
├── post-service/           # Property listings
├── favourite-service/      # User favorites
├── logger-service/         # Centralized logging
├── mail-service/           # Email notifications
├── listener-service/       # Message queue consumer
├── front-end/              # Web UI (Go templates)
├── project/                # Windows Docker Compose
├── linux_project/          # Linux Docker Compose
└── public/                 # Static assets
```

## Message Queue Events

### Log Events (logs_topic exchange)

- `log.INFO` - Information logs
- `log.WARNING` - Warning logs
- `log.ERROR` - Error logs

### App Events (app_events exchange)

- `mail.send` - Send generic email
- `mail.verification` - Send verification email
- `mail.password_reset` - Send password reset email
- `notification.send` - Send notification

## Makefile Commands

| Command               | Description                   |
| --------------------- | ----------------------------- |
| `make up`             | Start all containers          |
| `make up_build`       | Build and start containers    |
| `make down`           | Stop all containers           |
| `make build_broker`   | Build broker service binary   |
| `make build_auth`     | Build auth service binary     |
| `make build_logger`   | Build logger service binary   |
| `make build_mail`     | Build mail service binary     |
| `make build_listener` | Build listener service binary |

## License

MIT License
