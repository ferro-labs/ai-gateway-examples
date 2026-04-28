# Embedded Example

Shows how to embed the AI Gateway as an HTTP handler inside your own `net/http` server, mounting it alongside existing application routes.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./embedded
```

Then in another terminal:

```bash
# Your existing app route
curl http://localhost:8080/api/hello

# Gateway route
curl -X POST http://localhost:8080/ai/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hello"}]}'
```

## What it does

1. Creates a standard `net/http` mux with existing app routes (`/api/hello`, `/health`)
2. Mounts the gateway handler at `/ai/v1/chat/completions`
3. Serves everything from a single server on `:8080`

## Expected output

```
Provider: openai
Server listening on :8080
  App route:     GET  /api/hello
  Gateway route: POST /ai/v1/chat/completions
```
