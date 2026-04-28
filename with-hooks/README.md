# Event Hooks Example

Demonstrates gateway event hooks for telemetry and audit logging. Hooks fire asynchronously after each request without blocking the response path.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./with-hooks
```

## What it does

1. Creates a gateway and registers an event hook via `gw.AddHook(fn)`
2. Sends a successful request and an intentionally invalid request
3. The hook collects `gateway.request.completed` and `gateway.request.failed` events
4. Prints all received hook events with their payloads

## Expected output

```
Provider: openai  Model: gpt-4o-mini

--- Request 1: successful request ---
{
  "id": "chatcmpl-...",
  ...
}

--- Request 2: invalid request (empty messages) ---
Error (expected): messages must not be empty

--- Hook events received: 2 ---

Event 1:
{
  "subject": "gateway.request.completed",
  "data": { "model": "gpt-4o-mini", "latency_ms": 342, ... }
}

Event 2:
{
  "subject": "gateway.request.failed",
  "data": { "error": "messages must not be empty", ... }
}
```
