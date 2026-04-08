# Cafeteria Delivery API

A REST API for a cafeteria delivery platform built with Go, Fiber, and PostgreSQL following Hexagonal Architecture (
Ports & Adapters).

## Domain

The system manages a cafeteria menu and customer orders.

### Entities

**Item Category** — groups menu items (e.g. Mains, Drinks, Desserts, Snacks).

**Item** — a menu product with a price and inventory quantity. Belongs to a category.

**Order** — a customer order containing one or more items. Goes through the following status lifecycle:

| Status    | Value |
|-----------|-------|
| Pending   | 1     |
| Confirmed | 2     |
| Ready     | 3     |
| Delivered | 4     |
| Cancelled | 5     |

### Business Rules

- A new order is always created with status `Pending`.
- All items in an order must exist in the catalogue.
- Each ordered item's quantity must not exceed the available inventory.
- Unit prices are captured at the time of order creation.
- Total price is calculated as the sum of `unit_price × quantity` for all items.

---

## Running Locally

### Prerequisites

- Go 1.25+
- PostgreSQL running locally
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI installed

### Steps

```bash
# 1. Clone and install dependencies
git clone https://github.com/SE510-DS-1-Microservices-26/building-microservice-system-253.git
cd cafeteria-delivery
go mod download

# 2. Copy and fill in environment variables
cp .env.example .env

# 3. Run migrations
migrate -path database/migrations \
  -database "postgres://dev:root1234@localhost:5432/cafeteria_delivery?sslmode=disable" up

# 4. Start the server
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`.

---

## Running with Docker

### Prerequisites

- Docker + Docker Compose

### Steps

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Build and start all services (API + DB + migrations + seed)
make build

# 3. Stop all services
make down

# 4. Reset database (wipe volumes and restart)
make db-reset

# 5. Re-seed the database
make seed
```

The API will be available at `http://localhost:8000`.

---

## API Examples

### Health Check

```bash
curl http://localhost:8000/health
```

---

### Item Categories

**List categories**

```bash
curl http://localhost:8000/item-categories
```

**Get a category**

```bash
curl http://localhost:8000/item-categories/1
```

**Create a category**

```bash
curl -X POST http://localhost:8000/item-categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Specials"}'
```

**Update a category**

```bash
curl -X PUT http://localhost:8000/item-categories/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Daily Specials"}'
```

**Delete a category**

```bash
curl -X DELETE http://localhost:8000/item-categories/1
```

---

### Items

**List items**

```bash
curl http://localhost:8000/items
```

**List items filtered by category**

```bash
curl "http://localhost:8000/items?category_id=1&limit=5&offset=0"
```

**Get an item**

```bash
curl http://localhost:8000/items/1
```

**Create an item**

```bash
curl -X POST http://localhost:8000/items \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Grilled Chicken",
    "description": "Served with seasonal vegetables",
    "price": 12.99,
    "quantity": 20
  }'
```

**Update an item**

```bash
curl -X PUT http://localhost:8000/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grilled Chicken",
    "description": "Served with seasonal vegetables",
    "price": 13.99,
    "quantity": 15
  }'
```

**Delete an item**

```bash
curl -X DELETE http://localhost:8000/items/1
```

---

### Orders

**List orders**

```bash
curl http://localhost:8000/orders
```

**List orders filtered by user and status**

```bash
curl "http://localhost:8000/orders?user_id=1&status=1"
```

**Get an order**

```bash
curl http://localhost:8000/orders/1
```

**Create an order**

```bash
curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"item_id": 1, "quantity": 2},
      {"item_id": 4, "quantity": 1}
    ]
  }'
```

**Update order status**

```bash
curl -X PUT http://localhost:8000/orders/1 \
  -H "Content-Type: application/json" \
  -d '{"status": 2}'
```

**Delete an order**

```bash
curl -X DELETE http://localhost:8000/orders/1
```

---

## Testing

```bash
# Run all tests
make test

# Unit tests only (services)
make unit-test

# API tests only (handlers)
make api-test
```

---

## Swagger Docs

```bash
# Generate docs
make swagger-gen

# View in browser (server must be running)
open http://localhost:8000/docs/
```