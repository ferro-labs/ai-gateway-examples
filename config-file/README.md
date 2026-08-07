# Config File Example

Loads the entire gateway configuration from [`gateway.yaml`](./gateway.yaml) with `config.LoadConfig` instead of building `config.Config` in Go. This is how a real deployment is configured — edit YAML, restart, no recompile.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./config-file
```

The YAML targets the `openai` provider, so this example needs `OPENAI_API_KEY`. To use a different provider, change the target's `virtual_key` in `gateway.yaml` and register that provider in `main.go`.

## What it does

1. Loads `gateway.yaml` (strategy, targets, and a `fast` alias) via `config.LoadConfig`
2. Creates the gateway from the loaded config and registers the OpenAI provider
3. Routes a request using the **alias** `fast` — the gateway resolves it to `gpt-4o-mini`

## Expected output

```
Loaded config: strategy=single, targets=1, aliases=map[fast:gpt-4o-mini]

Routing request with model alias "fast"...
Alias "fast" resolved to model: gpt-4o-mini-2024-07-18
{ ... }
```

Config decoding is strict: a misspelled key in `gateway.yaml` is rejected at load time rather than silently ignored.
