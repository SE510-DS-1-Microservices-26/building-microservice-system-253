# Practice 7 — Workflow / Saga + Kubernetes Deployment

## Overview

For this practice we extended our existing Cafeteria Delivery microservice system by adding a Workflow Service that implements the Saga pattern, and prepared the entire system for Kubernetes deployment. The goal was to orchestrate a multi-step business process across services with a proper compensation (rollback) mechanism when something goes wrong.

---

## What We Built

### Workflow Service (Saga / Process Manager)

We introduced a brand new `workflow-service` written in Go using the same Fiber framework and Clean Architecture approach as the rest of the system. The service manages the `place-order` workflow, which coordinates three steps across existing services:

1. **Validate the user** — calls the Users Service to confirm the user actually exists before doing anything else
2. **Create the order** — calls the Core Service to create the order with the provided items
3. **Complete** — marks the workflow as successfully finished

Each step updates the workflow state in the database, so you can always query the current status and see exactly where the process is.

#### State Machine

One thing we're particularly happy about is the state machine approach in the domain layer. Instead of just having a string field that anyone can set freely, we defined explicit allowed transitions:

```
started → user_validated → order_created → completed
started → compensating → compensated
any step → failed
```

The `Transition()` method on the entity validates every state change and rejects invalid ones. This makes the saga logic easy to reason about and hard to misuse.

#### Compensation Path

If the order creation step fails for any reason (item out of stock, invalid data, Core Service down), the workflow automatically switches to a `compensating` state and calls the Core Service to cancel any order that may have been partially created. The error reason is saved in the `last_error` field so nothing gets lost silently.

#### Database

The `workflow_instances` table stores the full state of every workflow run:

| Column | Description |
|---|---|
| `workflow_id` | UUID v7 primary key (time-sortable) |
| `type` | Workflow type, e.g. `place-order` |
| `state` | Current state in the saga |
| `payload` | Original request as JSONB |
| `last_error` | Populated only when something fails |
| `created_at` / `updated_at` | Standard timestamps |

We chose `uuid_generate_v7()` over the standard `gen_random_uuid()` because UUIDv7 embeds a timestamp, which makes them naturally sortable by creation time — a nice win for queries and debugging.

#### API

The service exposes two endpoints:

```
POST /workflows/place-order   — starts the saga
GET  /workflows/:workflowId   — returns current state and any error
```

Both are routed through the existing API Gateway, which we extended with a `/workflow/*` proxy route.

---

### Kubernetes Manifests

We created a `/k8s` folder with manifests for every component in the system. The structure is straightforward — one YAML file per service, each containing everything that service needs.

**Services covered:**
- `core.yaml` — Core Service + PostgreSQL StatefulSet
- `users.yaml` — Users Service + PostgreSQL StatefulSet
- `notification.yaml` — Notification Service + PostgreSQL StatefulSet
- `workflow.yaml` — Workflow Service + PostgreSQL StatefulSet
- `gateway.yaml` — API Gateway exposed via NodePort on port `30000`
- `rabbitmq.yaml` — RabbitMQ with management UI
- `namespace.yaml` — dedicated `cafeteria` namespace

Every Deployment follows the same pattern:
- Environment variables loaded from a **ConfigMap** (non-sensitive config) and a **Secret** (database credentials)
- **Readiness probe** — HTTP `/health` endpoint or `pg_isready` for databases — so Kubernetes only sends traffic when the service is actually ready
- **Liveness probe** — same endpoints, restarts the pod if it stops responding
- **Resource requests and limits** — small but realistic values (100m CPU / 128Mi memory for services, slightly more for databases) to prevent any single pod from starving others

Databases use **StatefulSets** with **PersistentVolumeClaims** so data survives pod restarts.

To deploy everything:

```bash
kubectl apply -f k8s/
kubectl get pods -n cafeteria
```

---

## Bug Fixes Along the Way

While building this we also fixed a few pre-existing issues in the codebase:

- **Swagger docs were never generated** — ran `swag init` for both `core` and `users` services and committed the generated `docs/` folder
- **`order_test.go` was calling `NewOrderService` with 3 arguments** instead of 4 — updated the test call
- **`order.go` would panic if publisher was nil** — added a nil guard before calling `PublishOrderCreated`, which also makes unit testing without a real RabbitMQ connection possible

---

## How to Verify

**Happy path** — start a valid workflow:
```bash
curl -X POST http://localhost:8000/workflow/workflows/place-order \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "items": [{"item_id": 1, "quantity": 2}]}'
```

Check the result:
```bash
curl http://localhost:8000/workflow/workflows/<workflow_id>
# state should be "completed"
```

**Compensation path** — trigger a failure by using a non-existent user:
```bash
curl -X POST http://localhost:8000/workflow/workflows/place-order \
  -H "Content-Type: application/json" \
  -d '{"user_id": 99999, "items": [{"item_id": 1, "quantity": 1}]}'
```

Check again:
```bash
curl http://localhost:8000/workflow/workflows/<workflow_id>
# state should be "failed", last_error will explain why
```

---

## Pair Programming Sessions

| Session | Focus |
|---|---|
| Session 1 | Designed the saga flow and state machine together, agreed on domain model |
| Session 2 | Workflow service implementation + Kubernetes manifests review |
