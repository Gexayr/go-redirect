# Platform Architecture and Implementation

## Overview
This document describes the architecture and implementation of a backend platform built with Golang, RabbitMQ, and MySQL. The platform provides client authentication, request redirection management, and request logging capabilities. It is containerized using Docker Compose and operates entirely as an API.

## Architecture

### Components:
1. **API Gateway (Golang)**
   - Handles HTTP requests from clients
   - Manages client authentication and authorization
   - Processes redirect requests
   - Sends request data to RabbitMQ for logging
   - Implements JWT-based authentication

2. **Message Queue (RabbitMQ)**
   - Acts as a broker for asynchronous request logging
   - Ensures reliable message delivery to the database worker
   - Provides message persistence and retry mechanisms

3. **Database Worker (Golang)**
   - Consumes messages from RabbitMQ
   - Processes and stores request data in MySQL
   - Handles redirect history tracking
   - Implements error handling and logging

4. **Database (MySQL)**
   - Stores client information and credentials
   - Maintains redirect mappings and history
   - Logs request data for analysis
   - Implements proper indexing for performance

5. **Docker Compose**
   - Orchestrates all services
   - Manages service dependencies
   - Provides isolated development environment

## Implementation Details

### API Gateway
- Implements RESTful API endpoints
- Handles client registration and authentication
- Manages redirect mapping creation and retrieval
- Processes redirect requests with hash-based routing
- Implements JWT token generation and validation
- Uses middleware for request validation and authentication

### Authentication System
- JWT-based authentication
- Secure password hashing with bcrypt
- Token expiration (24 hours)
- Protected API endpoints
- Client-specific redirect management

### RabbitMQ Integration
- Asynchronous message processing
- Message persistence
- Error handling and retries
- Queue management
- Message acknowledgment system

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
│   │   │   ├── client.go
│   │   │   └── redirect.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   └── logging.go
│   │   └── router.go
│   ├── auth/
│   │   └── jwt.go
│   ├── models/
│   │   ├── client.go
│   │   ├── redirect.go
│   │   └── request.go
│   ├── repository/
│   │   ├── mysql/
│   │   │   ├── client_repository.go
│   │   │   ├── redirect_repository.go
│   │   │   └── request_repository.go
│   │   └── rabbitmq/
│   │       ├── publisher.go
│   │       └── consumer.go
│   └── service/
│       ├── client_service.go
│       └── redirect_service.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
├── scripts/
│   └── migrations/
│       └── 000_init.sql
├── docker-compose.yml
├── env.example
└── go.mod
```

## Workflow

### Client Registration and Authentication
1. Client registers with username, password, and email
2. System validates input and checks for duplicates
3. Password is hashed and client record is created
4. JWT token is generated and returned to client

### Redirect Management
1. Authenticated client creates redirect mapping
2. System generates unique 6-symbol hash
3. Mapping is stored in database
4. Client can retrieve their mappings

### Request Processing
1. Request arrives at API Gateway
2. System extracts hash from URL
3. If hash exists:
   - Retrieves redirect mapping
   - Performs redirection
   - Logs request data
4. If no hash:
   - Logs request data only

### Message Processing
1. API Gateway publishes message to RabbitMQ
2. Database worker consumes message
3. Worker processes message:
   - Saves request to request_logs
   - Updates redirect_history if applicable
4. Message is acknowledged after successful processing
5. Failed messages are requeued for retry

## Security Considerations

- JWT tokens for authentication
- Password hashing with bcrypt
- Protected API endpoints
- Input validation
- Rate limiting
- Secure headers
- Environment variable configuration

## Configuration

The platform uses environment variables for configuration. See `env.example` for all required variables and their descriptions. Key configuration areas include:

- API Gateway settings
- Database connection details
- RabbitMQ configuration
- JWT settings
- Logging configuration
- Service URLs

## Conclusion

This platform provides a robust solution for request redirection and logging with:
- Secure client authentication
- Efficient request processing
- Reliable message handling
- Scalable architecture
- Easy deployment with Docker
- Comprehensive logging and monitoring 