# Expense Tracker Service

A simple Go REST API for tracking expenses with PostgreSQL database.

## API Endpoints

- `POST /expenses` - Create a new expense
- `GET /expenses` - List all expenses
- `GET /expenses?id={id}` - Get a specific expense
- `PUT /expenses?id={id}` - Update an expense
- `DELETE /expenses?id={id}` - Delete an expense

## Prerequisites

- Go 1.21+
- Docker and Docker Compose

## Setup

1. Start PostgreSQL:
```bash
docker-compose up -d
```

2. Run the service:
```bash
go run main.go
```

The service will start on port 8002.

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (default: localhost postgres)

## Example Usage

Create an expense:
```bash
curl -X POST http://localhost:8002/expenses \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.00, "category": "Food", "description": "Lunch", "date": "2026-05-16T12:00:00Z"}'
```

List expenses:
```bash
curl http://localhost:8002/expenses
```
