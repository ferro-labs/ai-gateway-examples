# Guardrails Example

Demonstrates using guardrail plugins to validate and filter requests before they reach the provider. Includes a word-filter and a max-token limit.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./with-guardrails
```

## What it does

1. Registers `word-filter` and `max-token` guardrail plugins
2. Configures blocked words: "password", "secret", "api_key"
3. Configures token limits: max 4096 tokens, max 50 messages
4. Sends a clean request (passes through) and a blocked request (rejected)

## Expected output

```
Provider: openai  Model: gpt-4o-mini

--- Request 1: clean message (should pass) ---
{
  "id": "chatcmpl-...",
  ...
}

--- Request 2: blocked word 'password' (should be rejected) ---
Rejected: blocked word detected: password
```
