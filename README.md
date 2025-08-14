# Auth Service

A simple authentication service written in Go using Gin with PostgreSQL database integration.

## Features
- User registration with password hashing
- User login with credential validation
- PostgreSQL database integration with GORM
- Configurable via YAML and environment variables
- Docker ready
- Health check endpoint

## Prerequisites
- Go 1.24+
- Docker and Docker Compose (for local development)
- PostgreSQL (or use Docker Compose)

## Database Setup

### Using Docker Compose (Recommended)
```bash
# Start PostgreSQL
docker-compose up -d postgres

# Wait a few seconds for the database to be ready
sleep 5
```

### Manual PostgreSQL Setup
1. Install PostgreSQL
2. Create a database named `auth_service`
3. Update the database configuration in `configs/config.yaml`

## Environment Variables

Create a `.env` file or set these environment variables:

```bash
# Database credentials (REQUIRED)
export DB_USER=tranquangkha
export DB_PASSWORD=

# Optional: Override other settings
export APP_NAME=auth-service
export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=ck_auth
export DB_SSL_MODE=disable
```

## Run locally
```bash
# Set environment variables
export DB_USER=tranquangkha
export DB_PASSWORD=

# Start the database (if using Docker Compose)
docker-compose up -d postgres

# Run the application
go run ./cmd/server
```

## API Endpoints

- `GET /health` - Health check
- `POST /register` - User registration
- `POST /login` - User login

### Example Usage

Register a new user:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

Login:
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

Health check:
```bash
curl http://localhost:8080/health
```
