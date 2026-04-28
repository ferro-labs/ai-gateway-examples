# Fallback Example

Demonstrates the fallback routing strategy. The gateway tries providers in order — if the primary fails, it retries then automatically falls back to the next.

## Run

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./fallback
```

## What it does

1. Detects all configured providers
2. Builds a fallback chain with retry (3 attempts) on the primary provider
3. Routes a request — on failure, falls back through providers in order
4. Prints the response from whichever provider succeeded

## Expected output

```
Registered: openai
Registered: anthropic

Routing request (model=gpt-4o-mini, strategy=fallback)...
{
  "id": "chatcmpl-...",
  ...
}
```
