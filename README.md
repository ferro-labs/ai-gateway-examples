# Ferro Labs Examples

A central repository of runnable examples for all Ferro Labs projects, organized by product.

## Projects

| Folder | Project | Language | Description |
|--------|---------|----------|-------------|
| [`ferrotunnel-examples/`](./ferrotunnel-examples/) | [FerroTunnel](https://github.com/ferro-labs/ferrotunnel) | Rust | Secure, embedded, API-first tunneling with public URLs |
| [`ai-gateway-examples/`](./ai-gateway-examples/) | [AI Gateway](https://github.com/ferro-labs/ai-gateway) | Go | Unified API for OpenAI, Anthropic, and 100+ LLMs |

---

## FerroTunnel Examples

→ **[ferrotunnel-examples/](./ferrotunnel-examples/)**

Examples covering embedded server/client setup, plugin development, TLS configuration, multi-tunnel routing, observability, and common real-world scenarios.

**Quick start:**
```bash
cargo run -p ferrotunnel-examples --example embedded_server
cargo run -p ferrotunnel-examples --example embedded_client
```

| Category | Examples |
|----------|---------|
| Basic | `embedded_server`, `embedded_client`, `auto_reconnect` |
| Plugins | `custom_plugin`, `header_filter`, `ip_blocklist`, `plugin_chain` |
| Advanced | `tls_config`, `multi_tunnel`, `http2_connection_pooling`, `custom_pool_config` |
| Operational | `server_graceful_shutdown`, `server_observability` |
| Scenarios | `expose_local_dev`, `receive_webhooks_locally`, `websocket_tunnel` |

---

## AI Gateway Examples

→ **[ai-gateway-examples/](./ai-gateway-examples/)**

Examples covering direct provider calls, custom plugins, embedded HTTP handler, fallback and load-balance routing, circuit breakers, guardrails, event hooks, and MCP tool-call integration.

**Quick start:**
```bash
export OPENAI_API_KEY=sk-...
cd ai-gateway-examples && go run ./basic
```

| Example | Description |
|---------|-------------|
| `basic` | Send a request to the first configured provider |
| `custom-plugin` | Write and register a custom plugin |
| `embedded` | Mount the gateway inside your own HTTP server |
| `fallback` | Automatic provider fallback on failure |
| `loadbalance` | Weighted load balancing across providers |
| `with-circuit-breaker` | Per-provider circuit breaker |
| `with-guardrails` | Word-filter and token-limit guardrails |
| `with-hooks` | Event hooks for telemetry and audit logging |
| `with-mcp` | MCP (Model Context Protocol) tool-call integration |
