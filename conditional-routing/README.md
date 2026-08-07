# Conditional Routing Example

Routes different models to different providers through one gateway endpoint, using the `conditional` strategy. This is the building block for cost/tier routing — point a cheap model at a cheap provider and a premium model at a premium one, and let callers pick by model name without knowing which provider serves it.

## Run

Requires **two** provider keys:

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./conditional-routing
```

## What it does

1. Registers the first two configured providers as targets
2. Configures one condition per provider: `model == <that provider's model> → route to that provider`
3. Sends each model and prints which provider actually served it

Unmatched models fall through to the default (the first target).

## Expected output

```
Routing rules:
  model "gpt-4o-mini" -> openai
  model "claude-3-5-haiku-latest" -> anthropic

  model=gpt-4o-mini                             -> served by "openai"
  model=claude-3-5-haiku-latest                 -> served by "anthropic"
```

The conditions key on the request's `model` field. To route on prompt content instead (e.g. send prompts mentioning "code" to a stronger model), use `mode: content-based` with `content_conditions`.
