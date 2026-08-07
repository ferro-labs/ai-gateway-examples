# Guardrails Example

Demonstrates enabling the gateway's built-in guardrail plugins — `word-filter` and `max-token` — through configuration, so requests are validated and filtered before they reach the provider. To write your own plugin instead, see [custom-plugin](../custom-plugin).

## Run

```bash
OPENAI_API_KEY=sk-... go run ./with-guardrails
```

## What it does

1. Enables the built-in `word-filter` and `max-token` guardrail plugins via config
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
Rejected: request blocked by content policy
```

The blocked word is logged server-side but deliberately kept out of the
client-facing message, so the blocklist can't be probed one word at a time.
