# Cafeteria Delivery — Async Messaging

Extends the microservices platform with asynchronous event publishing via RabbitMQ. When an order is created in the Core
Service, an integration event is published to RabbitMQ and consumed by the new Notification Service, which stores a
persistent notification record.

## Architecture

| Service                  | Responsibility                                                                  |
|--------------------------|---------------------------------------------------------------------------------|
| **Core Service**         | Menu items, categories, orders. Publishes `core-item.created` on order creation |
| **Users Service**        | User CRUD. Validates users for Core on order creation                           |
| **Notification Service** | Consumes `core-item.created` events, stores notification records                |
| **Gateway**              | Dumb reverse proxy — routes `/core/*` and `/users/*`                            |
| **RabbitMQ**             | Message broker                                                                  |

---

## Messaging Contract

### Exchange

| Property | Value       |
|----------|-------------|
| Name     | `cafeteria` |
| Type     | `topic`     |
| Durable  | yes         |

### Queue

| Property    | Value                            |
|-------------|----------------------------------|
| Name        | `notification.core-item.created` |
| Durable     | yes                              |
| Routing key | `core-item.created`              |

### Event: `core-item.created`

Published by Core Service on every successful order creation.

**Payload**

```json
{
  "event_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "occurred_at": "2025-04-09T14:32:10.123456Z",
  "correlation_id": "f0e1d2c3-b4a5-6789-dcba-098765432100",
  "core_item_id": 7,
  "owner_user_id": 1,
  "summary": "Order #7 placed by user 1"
}
```

| Field            | Type          | Description                                    |
|------------------|---------------|------------------------------------------------|
| `event_id`       | UUID          | Unique event identifier — used for idempotency |
| `occurred_at`    | UTC timestamp | When the order was placed                      |
| `correlation_id` | UUID          | Trace ID for the originating request           |
| `core_item_id`   | uint          | Order ID in the Core Service                   |
| `owner_user_id`  | uint \| null  | User who placed the order (null if anonymous)  |
| `summary`        | string        | Human-readable description                     |

### Idempotency

The `notifications` table has a `UNIQUE` constraint on `event_id`. Duplicate deliveries are silently ignored (
`ON CONFLICT (event_id) DO NOTHING`) and the message is acknowledged — no duplicates, no crashes on retry.

---

## Running with Docker

```bash
# Copy environment file
cp .env.example .env

# Build and start everything
make build

# Stop all services
make down

# Reset all databases (wipe volumes)
make db-reset
```

### Seed databases

```bash
make core-seed    # items, categories, orders
make users-seed   # 4 sample users (Alice, Bob, Carol, David)
```

---

## Service URLs

| Resource               | URL                                      |
|------------------------|------------------------------------------|
| Core API               | `http://localhost:8000/core/...`         |
| Users API              | `http://localhost:8000/users/...`        |
| Core Swagger           | `http://localhost:8000/core/docs/`       |
| Users Swagger          | `http://localhost:8000/users/docs/`      |
| RabbitMQ Management UI | `http://localhost:15672` (guest / guest) |

---

## Verifying Messages

### 1. Create an order to trigger an event

```bash
curl -X POST http://localhost:8000/core/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [{"item_id": 1, "quantity": 2}]
  }'
```

### 2. Check the notification was stored

```bash
docker compose exec notification_database \
  psql -U dev -d notification_db -c "SELECT * FROM notifications;"
```

Expected output:

```
 id | event_id | core_item_id | owner_user_id |          summary          |       occurred_at        |         created_at
----+----------+--------------+---------------+---------------------------+--------------------------+----------------------------
  1 | a1b2c3.. |            7 |             1 | Order #7 placed by user 1 | 2025-04-09 14:32:10.123  | 2025-04-09 14:32:10.456
```

### 3. Verify via RabbitMQ Management UI

1. Open `http://localhost:15672` → log in with `guest` / `guest`
2. Go to **Queues and Streams** tab
3. Click `notification.core-item.created`
4. Scroll to **Get messages** → click **Get Message(s)** to peek at queued payloads
5. **Overview** tab shows real-time message rates while creating orders

### 4. Test idempotency — replay the same event

Send the same order twice — each gets a new `event_id` so both are stored. To test true idempotency (same event
replayed), check notification service logs — duplicate `event_id` is logged and silently skipped:

```bash
docker compose logs notification_service | grep "duplicate event"
```

---

## Troubleshooting

**Notification service fails to start — "failed to connect to rabbitmq"**
RabbitMQ has a healthcheck; the notification service waits for it. If it still fails, RabbitMQ may be slow to
initialize. Run:

```bash
docker compose restart notification_service
```

**No rows in `notifications` after creating an order**
Check Core Service logs — if publish fails it logs a warning but does not fail the order:

```bash
docker compose logs core_service | grep "failed to publish"
```

If you see this, RabbitMQ may not be healthy. Check: `docker compose ps rabbitmq`.

**`notification.core-item.created` queue does not appear in the Management UI**
The queue is declared on consumer startup. Check notification service is running:

```bash
docker compose ps notification_service
docker compose logs notification_service
```

**Order creation returns 500 with "failed to validate user"**
Users Service is down or the `user_id` does not exist. Either omit `user_id` (anonymous order) or check:

```bash
docker compose ps users_service
curl http://localhost:8000/users/1
```

---

## Testing

```bash
make test              # all tests
make unit-test         # unit tests for both services
make api-test          # handler tests for both services
```