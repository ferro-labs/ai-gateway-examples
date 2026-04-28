# MCP (Model Context Protocol) Example

Demonstrates the AI Gateway's MCP integration. The gateway discovers tools from connected MCP servers, injects them into LLM requests, and drives the agentic tool-call loop automatically.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./with-mcp
```

Works best with models that support tool/function calling (gpt-4o, claude-3-5-sonnet, llama-3.3-70b-versatile, etc.).

## What it does

1. Starts a local mock MCP server exposing a `get_weather` tool
2. Configures the gateway to connect to the MCP server
3. Sends a weather question to the LLM
4. The gateway automatically:
   - Injects the `get_weather` tool definition into the request
   - Resolves `tool_calls` responses via the MCP server
   - Re-routes with tool results appended
5. Prints the final LLM answer

## Expected output

```
Mock MCP server listening at http://127.0.0.1:54321/mcp
Provider: openai  Model: gpt-4o-mini

Sending request (agentic MCP loop active)...
[mock MCP] tools/call get_weather({"location":"San Francisco, CA"}) → {"location":"San Francisco, CA","temperature":62,...}
{
  "id": "chatcmpl-...",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "The current weather in San Francisco is 62°F (Partly Cloudy) with 72% humidity and 12 mph winds."
      }
    }
  ],
  ...
}
```
