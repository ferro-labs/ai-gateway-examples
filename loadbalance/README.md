# Load Balance Example

Demonstrates weighted load balancing across multiple providers. Requests are distributed probabilistically by weight.

## Run

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./loadbalance
```

Requires at least two provider keys.

## What it does

1. Detects all configured providers
2. Assigns weights: first provider gets 70%, rest split the remaining 30%
3. Sends 5 requests, showing which provider handled each one
4. Prints the provider name, model, and reply for each request

## Expected output

```
Load balancing across:
  openai (weight=70)
  anthropic (weight=30)

Sending 5 requests (model=gpt-4o-mini)...

  Request 1: {"id":"chatcmpl-...","provider":"openai","model":"gpt-4o-mini","reply":"request 1 ok"}
  Request 2: {"id":"chatcmpl-...","provider":"openai","model":"gpt-4o-mini","reply":"request 2 ok"}
  Request 3: {"id":"chatcmpl-...","provider":"anthropic","model":"gpt-4o-mini","reply":"request 3 ok"}
  ...
```
