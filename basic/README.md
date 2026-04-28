# Basic Example

Sends a single chat completion request to the first configured provider. The simplest possible usage of the AI Gateway library.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./basic
```

## What it does

1. Detects the first provider with a configured API key
2. Picks the first supported model from that provider
3. Sends a chat completion request
4. Prints the JSON response

## Expected output

```
Provider: openai  Model: gpt-4o-mini
{
  "id": "chatcmpl-...",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Why do programmers prefer dark mode? Because light attracts bugs!"
      }
    }
  ],
  ...
}
```
