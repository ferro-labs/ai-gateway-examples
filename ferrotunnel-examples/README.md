# FerroTunnel Examples

Examples demonstrating how to use [FerroTunnel](https://github.com/ferro-labs/ferrotunnel) — a secure, embedded, API-first tunneling library written in Rust — in your own applications.

## Directory Structure

```
ferrotunnel-examples/
├── basic/                    # Getting started
│   ├── embedded_server.rs    # Embed a tunnel server in your app
│   ├── embedded_client.rs    # Embed a tunnel client in your app
│   └── auto_reconnect.rs     # Client with auto-reconnect
├── plugins/                  # Plugin development
│   ├── custom_plugin.rs      # Request counting and path blocking
│   ├── header_filter.rs      # Filter/modify HTTP headers
│   ├── ip_blocklist.rs       # Block requests by IP address
│   └── plugin_chain.rs       # Multiple plugins working together
├── advanced/                 # Advanced features
│   ├── tls_config.rs         # Configure TLS for secure connections
│   ├── multi_tunnel.rs       # Run multiple tunnels simultaneously
│   ├── http2_connection_pooling.rs  # HTTP/2 support and connection pooling
│   └── custom_pool_config.rs # Custom connection pool configuration
├── operational/              # Server lifecycle and observability
│   ├── server_graceful_shutdown.rs  # Graceful shutdown on SIGTERM/SIGINT
│   └── server_observability.rs     # Prometheus metrics and logging
└── scenarios/                # Common usage scenarios
    ├── expose_local_dev.rs         # Expose local dev server to the internet
    ├── receive_webhooks_locally.rs # Forward webhooks to local machine
    └── websocket_tunnel.rs         # WebSocket tunneling
```

## Quick Start

### Prerequisites

- Rust (stable) — [install](https://rustup.rs)
- A running FerroTunnel server (or use the embedded server example)

### Run an Example

```bash
# Start a local HTTP server to tunnel
python3 -m http.server 8000

# Terminal 1: Start the embedded server
cargo run -p ferrotunnel-examples --example embedded_server

# Terminal 2: Start the embedded client
cargo run -p ferrotunnel-examples --example embedded_client
```

## Examples

### Basic

| Example | Description |
|---------|-------------|
| `embedded_server` | Embed a full tunnel server inside your Rust application |
| `embedded_client` | Embed a tunnel client to expose a local service |
| `auto_reconnect` | Client that automatically reconnects on connection loss |

```bash
cargo run -p ferrotunnel-examples --example embedded_server
cargo run -p ferrotunnel-examples --example embedded_client
cargo run -p ferrotunnel-examples --example auto_reconnect
```

### Plugins

| Example | Description |
|---------|-------------|
| `custom_plugin` | Build a plugin that counts requests and blocks paths |
| `header_filter` | Remove sensitive headers and add security headers |
| `ip_blocklist` | Reject requests from blocked IP addresses |
| `plugin_chain` | Compose TokenAuth + RateLimit + Logger plugins |

```bash
cargo run -p ferrotunnel-examples --example custom_plugin
cargo run -p ferrotunnel-examples --example header_filter
cargo run -p ferrotunnel-examples --example ip_blocklist
cargo run -p ferrotunnel-examples --example plugin_chain
```

### Advanced

| Example | Description |
|---------|-------------|
| `tls_config` | Configure mutual TLS for the tunnel control plane |
| `multi_tunnel` | Run multiple tunnel clients for different local services |
| `http2_connection_pooling` | HTTP/2 auto-detection and connection pooling |
| `custom_pool_config` | Fine-tune the connection pool for your workload |

```bash
cargo run -p ferrotunnel-examples --example tls_config -- --mode server
cargo run -p ferrotunnel-examples --example multi_tunnel
cargo run -p ferrotunnel-examples --example http2_connection_pooling -- server
cargo run -p ferrotunnel-examples --example custom_pool_config
```

### Operational

| Example | Description |
|---------|-------------|
| `server_graceful_shutdown` | Handle SIGTERM/SIGINT and drain in-flight connections |
| `server_observability` | Prometheus metrics endpoint and structured logging |

```bash
cargo run -p ferrotunnel-examples --example server_graceful_shutdown
cargo run -p ferrotunnel-examples --example server_observability
```

### Scenarios

| Example | Description |
|---------|-------------|
| `expose_local_dev` | Share your local React/Next.js dev server with teammates |
| `receive_webhooks_locally` | Forward GitHub/Stripe webhooks to your local handler |
| `websocket_tunnel` | Tunnel WebSocket connections through FerroTunnel |

```bash
cargo run -p ferrotunnel-examples --example expose_local_dev
cargo run -p ferrotunnel-examples --example receive_webhooks_locally
cargo run -p ferrotunnel-examples --example websocket_tunnel
```

## More Information

- [FerroTunnel repository](https://github.com/ferro-labs/ferrotunnel)
- [API documentation](https://docs.rs/ferrotunnel)
