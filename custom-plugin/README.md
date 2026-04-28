# Custom Plugin Example

Demonstrates writing and registering a custom `plugin.Plugin` implementation. The example creates a request-ID injector that stamps a trace ID onto every request via plugin metadata.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./custom-plugin
```

## What it does

1. Defines a `requestIDPlugin` implementing the `plugin.Plugin` interface
2. Registers the plugin at the `before_request` stage
3. Routes a request — the plugin fires automatically before the provider call
4. Prints the trace ID and the provider's response

## Expected output

```
Provider: openai  Model: gpt-4o-mini
Custom plugin registered. Routing request...
[request-id] trace_id=req-1714300000000000  model=gpt-4o-mini
{
  "id": "chatcmpl-...",
  ...
}
```
