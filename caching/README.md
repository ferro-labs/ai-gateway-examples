# Response Caching Example

Enables the built-in `response-cache` plugin so identical requests are served from memory instead of calling the provider — saving latency and token cost.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./caching
```

Works with any of the 8 supported provider keys.

## What it does

1. Enables `response-cache` at two stages — `before_request` (serve a hit, skip the provider) and `after_request` (store the response). Both entries carry **identical config**, which is what makes them share one cache.
2. Sends the same request twice, timing each call
3. Shows that the second call is near-instant and returns the **same response ID** — proof it came from the cache, not the provider

## Expected output

```
Provider: openai  Model: gpt-4o-mini

Call 1 (cache miss):  480.0 ms  id=chatcmpl-...
Call 2 (cache hit):     0.1 ms  id=chatcmpl-...

Same response ID both times → call 2 was served from cache (4800x faster).
```

`max_age` (TTL seconds) and `max_entries` (capacity) tune the cache. Put guardrail plugins *ahead* of `response-cache` in the chain so a cache hit still runs them.
