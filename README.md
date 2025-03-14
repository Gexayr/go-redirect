# Platform

A backend platform built with Golang, RabbitMQ, and MySQL for handling request redirections and logging.

## Features

- Client authentication with JWT tokens
- Secure password hashing
- Request redirection with unique 6-symbol hashes
- Request logging and tracking
- Blacklist URL support
- Asynchronous message processing with RabbitMQ
- MySQL database for data persistence

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# API Gateway
API_PORT=8080
API_HOST=0.0.0.0

# MySQL
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=platform_db
MYSQL_USER=platform_user
MYSQL_PASSWORD=platform_password

# RabbitMQ
RABBITMQ_DEFAULT_USER=platform_user
RABBITMQ_DEFAULT_PASS=platform_password

# JWT
JWT_SECRET=your_jwt_secret_key
```

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd platform
```

2. Start the services:
```bash
docker-compose up --build
```

The platform will be available at `http://localhost:8080`.

## API Endpoints

### Authentication

#### Register a new client
```http
POST /api/register
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}
```

### Redirect Management (requires JWT token)

#### Create a new redirect mapping
```http
POST /api/redirects
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "redirect_url": "https://example.com",
    "redirect_url_black": "https://blacklist-example.com"
}
```

#### Get all redirect mappings
```http
GET /api/redirects
Authorization: Bearer <jwt_token>
```

### Redirect Access

#### Access redirect with hash
```http
GET /{hash}
```

## Authentication

The platform uses JWT (JSON Web Token) for authentication. When a client registers or logs in, they receive a JWT token that must be included in the Authorization header for protected endpoints.

Example:
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

The token expires after 24 hours, requiring the client to log in again.

## Database Schema

### Tables

#### clients
- `id` (BIGINT, PRIMARY KEY)
- `username` (VARCHAR(50), UNIQUE)
- `password_hash` (VARCHAR(255))
- `email` (VARCHAR(255), UNIQUE)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

#### redirect_mappings
- `id` (BIGINT, PRIMARY KEY)
- `client_id` (BIGINT, FOREIGN KEY)
- `hash` (VARCHAR(6), UNIQUE)
- `redirect_url` (TEXT)
- `redirect_url_black` (TEXT)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

#### request_logs
- `id` (BIGINT, PRIMARY KEY)
- `timestamp` (DATETIME)
- `ip_address` (VARCHAR(45))
- `request_url` (TEXT)
- `request_method` (VARCHAR(10))
- `request_headers` (JSON)
- `processing_status` (ENUM)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

#### redirect_history
- `id` (BIGINT, PRIMARY KEY)
- `request_log_id` (BIGINT, FOREIGN KEY)
- `original_url` (TEXT)
- `redirect_url` (TEXT)
- `redirect_type` (ENUM)
- `redirect_status` (INT)
- `redirect_timestamp` (DATETIME)

## Development

### Project Structure

```
platform/
├── cmd/
│   ├── api-gateway/
│   │   └── main.go
│   └── database-worker/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── router.go
│   ├── auth/
│   │   └── jwt.go
│   ├── models/
│   ├── repository/
│   └── service/
├── pkg/
│   ├── logger/
│   └── validator/
├── scripts/
│   └── migrations/
├── docker-compose.yml
└── go.mod
```

### Service Ports

- API Gateway runs on port 8080 by default
- Database Worker processes messages from RabbitMQ
- MySQL database runs on port 3306
- RabbitMQ management interface available on port 15672

### Message Processing Flow

The platform follows a specific flow for processing requests:

1. The API Gateway publishes messages to RabbitMQ
2. The database worker consumes messages and processes them:
   - Saves the request to `request_logs`
   - Extracts the hash from the URL
   - If a hash is found, gets the redirect URL from `redirect_mappings`
   - If a redirect URL is found, saves the redirect to `redirect_history`
3. The message is acknowledged only after all processing is complete
4. If any error occurs, the message is rejected and requeued

This flow ensures reliable message processing and data persistence, with automatic retries for failed operations.

### Running Tests

```bash
go test ./...
```

## License

[Your License] 