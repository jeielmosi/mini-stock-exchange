# Mini Stock Exchange

A mini stock exchange system designed to receive and process orders similar to how B3 or NYSE would work, allowing brokers to submit orders on behalf of their customers and resolve trades.

## Overview

This application implements a stock exchange matching engine that:
- Receives bid (buy) and ask (sell) orders from brokers
- Matches orders based on price and time priority
- Executes trades when matching orders are found
- Tracks order status and trade history

## Features

- Order submission with broker ID, customer document, order type, validity, stock symbol, price, and quantity
- Price matching logic (buyer's max price ≥ seller's min price)
- Partial order fulfillment support
- Chronological order resolution for same-price matches
- RESTful API for broker interactions
- Docker containerization for easy deployment
- Comprehensive test suite

## API Endpoints

### Broker Management
- `POST /brokers` - Register a new broker
- `GET /brokers/{id}` - Retrieve broker information

### Order Management
- `POST /orders` - Submit a new order (bid or ask)
- `GET /orders/{id}` - Get order status and details

### Trade Information
- `GET /trades` - List executed trades
- `GET /trades/{id}` - Get specific trade details

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (for persistence)

### Running with Docker Compose

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### Running Locally

```bash
# Install dependencies
go mod download

# Run the application
go run ./cmd/api/main.go
```

## Project Structure

```
├── cmd/                 # Application entry points
│   ├── api/             # Main API server
│   └── debug/           # Debug utilities
├── internal/            # Private application code
│   ├── controller/      # HTTP handlers
│   ├── dto/             # Data transfer objects
│   ├── entity/          # Core domain entities
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   ├── domain-service/  # Domain-specific services (matching engine)
│   └── usecase/         # Application use cases
├── scripts/             # Utility scripts
└── Dockerfile           # Container definition
```

## Technology Stack

- **Language**: Go 1.25
- **Framework**: chi (HTTP router)
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Testing**: Go testing framework

## Testing

Run the test suite:

```bash
go test ./...
```

Specific test categories:
- Unit tests: `go test ./internal/...`
- Integration tests: `go test ./internal/integration/...`
- API tests: See `scripts/tests/` for Python test scripts

## Configuration

Configuration is managed through environment variables:
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `SERVER_PORT`: HTTP server port (default: 8080)
