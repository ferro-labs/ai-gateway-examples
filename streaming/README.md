# Streaming Example

Demonstrates streaming chat completions through the AI Gateway. Tokens are printed in real-time as they arrive instead of waiting for the full response.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./streaming
```

## What it does

1. Creates a gateway with a single provider
2. Sends a streaming request (`stream: true`) via `gw.RouteStream()`
3. Reads `StreamChunk` objects from the returned channel
4. Prints each delta content as it arrives
5. Reports total token usage if available

## Expected output

```
Provider: openai  Model: gpt-4o-mini
Streaming response:

Servers talking,
Packets lost in the ether—
Retry, timeout, hope.

Total tokens: 28
```
