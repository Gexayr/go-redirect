# Platform

A backend platform built with Golang, RabbitMQ, and MySQL for request processing and redirection.

## Features

- Request processing and redirection
- Asynchronous request logging via RabbitMQ
- MySQL database storage
- Containerized with Docker Compose

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- MySQL 8.0
- RabbitMQ 3.12

## Project Structure

```
platform/
├── cmd/                    # Main entry points
├── internal/              # Private application code
├── pkg/                   # Reusable packages
└── scripts/              # Utility scripts
```

## Setup

1. Clone the repository
2. Copy `config/config.yaml.example` to `config/config.yaml` and update the values
3. Run `go mod download` to install dependencies
4. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

## Development

- API Gateway runs on port 8080 by default
- Database Worker processes messages from RabbitMQ
- MySQL database runs on port 3306
- RabbitMQ management interface available on port 15672

## Message Processing Flow

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

## License

MIT 