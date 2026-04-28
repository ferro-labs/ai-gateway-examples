// Package main demonstrates sending a request directly to any configured LLM provider.
//
// Set at least one provider key and run:
//
// OPENAI_API_KEY=sk-...       go run ./basic
// ANTHROPIC_API_KEY=sk-ant-... go run ./basic
// GROQ_API_KEY=gsk_...        go run ./basic
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ferro-labs/ai-gateway/providers"

	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func main() {
	// Pick the first provider that has a key configured.
	provider := shared.FirstProvider()

	model := provider.SupportedModels()[0]
	req := providers.Request{
		Model: model,
		Messages: []providers.Message{
			{Role: "user", Content: "Hello, tell me a short joke about programming."},
		},
	}

	if err := req.Validate(); err != nil {
		log.Fatalf("Invalid request: %v", err)
	}

	fmt.Printf("Provider: %s  Model: %s\n", provider.Name(), model)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := provider.Complete(ctx, req)
	if err != nil {
		cancel()
		log.Printf("Request failed: %v", err)
		os.Exit(1) //nolint:gocritic
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))
}
