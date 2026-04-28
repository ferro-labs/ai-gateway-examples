# Contributing to AI Gateway Examples

Thank you for your interest in contributing to the Ferro Labs AI Gateway examples!

## Guidelines

This repository follows the same contributing guidelines as the main AI Gateway project. Please see the [AI Gateway CONTRIBUTING.md](https://github.com/ferro-labs/ai-gateway/blob/main/CONTRIBUTING.md) for:

- Code style and conventions
- Commit message format (Conventional Commits)
- Pull request process
- Branching strategy

## Adding a New Example

1. Create a new directory under the root (e.g., `with-feature/`)
2. Add a `main.go` with a package-level doc comment explaining the example
3. Include proper error handling and context timeouts
4. Use environment variables for API keys (never hardcode)
5. Update the root `README.md` with your example in the directory listing
6. Test with at least one provider before submitting

## Running Examples

```bash
export OPENAI_API_KEY=sk-your-key
cd basic/
go run main.go
```

## Questions?

Open a [GitHub Discussion](https://github.com/ferro-labs/ai-gateway/discussions) or reach out on [Discord](https://discord.gg/YYSKrgBXMz).
