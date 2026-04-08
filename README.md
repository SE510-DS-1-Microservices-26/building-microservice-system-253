# Cafeteria Delivery — Microservices

A cafeteria delivery platform split into three independent services communicating over HTTP.

## Architecture

| Service           | Responsibility                                   |
|-------------------|--------------------------------------------------|
| **Core Service**  | Menu items, categories, orders                   |
| **Users Service** | User CRUD                                        |
| **Gateway**       | Reverse proxy — routes requests to Core or Users |

All external traffic goes through the **Gateway** on port `8000`.

### Inter-service Communication

When creating an order with a `user_id`, Core calls `GET /users/{id}` on the Users Service to validate the user exists
before persisting the order. If Users is unreachable or returns a non-200 response, the order creation fails.

---

## Running with Docker

### Prerequisites

- Docker + Docker Compose

### Steps

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Build and start everything
make build

# 3. Stop all services
make down

# 4. Reset databases (wipe volumes)
make db-reset
```

### Re-seed databases

```bash
make core-seed    # re-seed core DB (items, categories, orders)
make users-seed   # re-seed users DB (4 sample users)
```

---

## Service URLs

All requests go through the Gateway:

| Resource                         | Base URL                            |
|----------------------------------|-------------------------------------|
| Core (items, categories, orders) | `http://localhost:8000/core/...`    |
| Users                            | `http://localhost:8000/users/...`   |
| Core Swagger UI                  | `http://localhost:8000/core/docs/`  |
| Users Swagger UI                 | `http://localhost:8000/users/docs/` |

---

## API Examples

### Users

**List users**

```bash
curl http://localhost:8000/users
```

**Get a user**

```bash
curl http://localhost:8000/users/1
```

**Create a user**

```bash
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Eve Adams", "email": "eve@example.com"}'
```

**Update a user**

```bash
curl -X PUT http://localhost:8000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice.new@example.com"}'
```

**Delete a user**

```bash
curl -X DELETE http://localhost:8000/users/1
```

---

### Item Categories

**List categories**

```bash
curl http://localhost:8000/core/item-categories
```

**Create a category**

```bash
curl -X POST http://localhost:8000/core/item-categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Specials"}'
```

---

### Items

**List items**

```bash
curl http://localhost:8000/core/items
```

**List items by category**

```bash
curl "http://localhost:8000/core/items?category_id=1&limit=5"
```

**Create an item**

```bash
curl -X POST http://localhost:8000/core/items \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Grilled Chicken",
    "description": "Served with seasonal vegetables",
    "price": 12.99,
    "quantity": 20
  }'
```

---

### Orders

**Create an order without a user**

```bash
curl -X POST http://localhost:8000/core/orders \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"item_id": 1, "quantity": 2},
      {"item_id": 4, "quantity": 1}
    ]
  }'
```

**Create an order tied to a user** (Core will validate the user exists)

```bash
curl -X POST http://localhost:8000/core/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"item_id": 1, "quantity": 1}
    ]
  }'
```

**Update order status**

```bash
curl -X PUT http://localhost:8000/core/orders/1 \
  -H "Content-Type: application/json" \
  -d '{"status": 2}'
```

**List orders by user and status**

```bash
curl "http://localhost:8000/core/orders?user_id=1&status=1"
```

---

## Order Status Reference

| Value | Meaning   |
|-------|-----------|
| 1     | Pending   |
| 2     | Confirmed |
| 3     | Ready     |
| 4     | Delivered |
| 5     | Cancelled |

---

## When Users Service is Down

Placing an order **without** `user_id` — works fine, Users Service is not called:

```bash
# Returns 201
curl -X POST http://localhost:8000/core/orders \
  -H "Content-Type: application/json" \
  -d '{"items": [{"item_id": 1, "quantity": 1}]}'
```

Placing an order **with** `user_id` — fails because Core cannot reach Users to validate:

```bash
# Returns 500 — "failed to validate user: failed to call users service: ..."
curl -X POST http://localhost:8000/core/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "items": [{"item_id": 1, "quantity": 1}]}'
```

All other Core endpoints (items, categories, orders without user_id) remain fully operational regardless of Users
Service availability.

---

## Testing

```bash
make test              # all tests
make unit-test         # unit tests for both services
make unit-test-core    # core service unit tests only
make unit-test-users   # users service unit tests only
make api-test          # API (handler) tests for both services
make api-test-core     # core handler tests only
make api-test-users    # users handler tests only
```

---

## Swagger Docs

```bash
# Regenerate docs
make swagger-gen

# View in browser (stack must be running)
open http://localhost:8000/core/docs/
open http://localhost:8000/users/docs/
```