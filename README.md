# AI Gateway Examples

Examples demonstrating how to use [Ferro Labs AI Gateway](https://github.com/ferro-labs/ai-gateway) — an open-source AI gateway written in Go with a unified API for OpenAI, Anthropic, Bedrock, Azure, and 100+ LLMs.

## Directory Structure

```text
.
├── basic/                    # Direct provider request — the simplest use case
├── streaming/                # Streaming chat completions (real-time token output)
├── custom-plugin/            # Write and register a custom plugin
├── embedded/                 # Embed the gateway inside an existing HTTP server
├── fallback/                 # Automatic fallback between providers on failure
├── loadbalance/              # Weighted load balancing across multiple providers
├── with-circuit-breaker/     # Per-provider circuit breaker configuration
├── with-guardrails/          # Built-in word-filter and token-limit guardrails
├── with-hooks/               # Event hooks for telemetry and audit logging
├── with-mcp/                 # MCP (Model Context Protocol) tool-call integration
└── shared/                   # Shared provider helpers (used by all examples)
```

## Quick Start

### Prerequisites

- Go 1.22+ — [install](https://go.dev/dl/)
- At least one LLM provider API key (OpenAI, Anthropic, Groq, Gemini, Mistral, Together, Cohere, or DeepSeek)

### Run an Example

```bash
# Set a provider API key
export OPENAI_API_KEY=sk-...
# or
export ANTHROPIC_API_KEY=sk-ant-...
# or any other supported provider key

# Run the basic example
go run ./basic
```

## Examples

### basic

Sends a single chat completion request to the first provider for which an API key is configured. Great starting point.

```bash
OPENAI_API_KEY=sk-... go run ./basic
```

### streaming

Streams chat completion tokens in real-time via `gw.RouteStream()`. Tokens are printed as they arrive instead of waiting for the full response.

```bash
OPENAI_API_KEY=sk-... go run ./streaming
```

### custom-plugin

Demonstrates writing a `plugin.Plugin` implementation (a request-ID injector) and registering it at the `before_request` stage.

```bash
OPENAI_API_KEY=sk-... go run ./custom-plugin
```

### embedded

Shows how to embed the gateway as an HTTP handler inside your own `net/http` server, mounting the AI route alongside your existing application routes.

```bash
OPENAI_API_KEY=sk-... go run ./embedded
# Then: curl -X POST http://localhost:8080/ai/v1/chat/completions -d '...'
```

### fallback

Configures `ModeFallback`: the gateway tries providers in order, retrying on the primary before switching to the next. Requires at least one provider key; works better with two.

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./fallback
```

### loadbalance

Configures `ModeLoadBalance` with weighted targets. The first provider receives 70% of requests, the rest share the remaining 30%. Requires at least two provider keys.

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./loadbalance
```

### with-circuit-breaker

Attaches a `CircuitBreakerConfig` to the primary provider target. After two consecutive failures the circuit opens for 10 s, routing all traffic to the fallback.

```bash
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./with-circuit-breaker
```

### with-guardrails

Loads the built-in `word-filter` and `max-token` guardrail plugins via `PluginConfig`. Requests containing blocked words are rejected before reaching the provider.

```bash
OPENAI_API_KEY=sk-... go run ./with-guardrails
```

### with-hooks

Registers an event hook via `gw.AddHook(fn)`. The hook receives `gateway.request.completed` and `gateway.request.failed` events asynchronously after each request — useful for telemetry, cost tracking, or audit logging.

```bash
OPENAI_API_KEY=sk-... go run ./with-hooks
```

### with-mcp

Starts a local mock MCP server, wires it into the gateway, and sends a tool-calling request. The gateway injects tool definitions, drives the agentic tool-call loop, and returns the final answer.

```bash
OPENAI_API_KEY=sk-... go run ./with-mcp
# Works best with models that support tool/function calling (gpt-4o, claude-3-5-sonnet, etc.)
```

## Supported Provider Environment Variables

| Variable | Provider |
|----------|----------|
| `OPENAI_API_KEY` | OpenAI |
| `ANTHROPIC_API_KEY` | Anthropic |
| `GROQ_API_KEY` | Groq |
| `GEMINI_API_KEY` | Google Gemini |
| `MISTRAL_API_KEY` | Mistral |
| `TOGETHER_API_KEY` | Together AI |
| `COHERE_API_KEY` | Cohere |
| `DEEPSEEK_API_KEY` | DeepSeek |

## More Information

- [AI Gateway repository](https://github.com/ferro-labs/ai-gateway)
- [Documentation](https://github.com/ferro-labs/ferrolabs-docs)
