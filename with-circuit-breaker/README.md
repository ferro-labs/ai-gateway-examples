# Circuit Breaker Example

Demonstrates per-provider circuit breaker configuration. After consecutive failures, the circuit opens and routes all traffic to the fallback provider.

## Run

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./with-circuit-breaker
```

## What it does

1. Sets up a two-provider fallback chain
2. Configures an aggressive circuit breaker on the primary (threshold: 2 failures, timeout: 10s)
3. Sends a request — routes to primary if circuit is closed
4. To observe the circuit opening: set the primary key to an invalid value

## Expected output

```
Registered: openai (models: 14)
Registered: anthropic (models: 8)

Primary provider : openai  (circuit breaker: threshold=2, timeout=10s)
Fallback provider: anthropic
Model            : gpt-4o-mini

--- Request 1: normal (circuit closed) ---
{
  "id": "chatcmpl-...",
  ...
}

Tip: set the primary key to an invalid value and rerun — after
     2 failures the circuit opens and the fallback takes over.
```
