// Package main demonstrates the fallback routing strategy.
//
// The gateway tries providers in order; if the first fails it automatically
// retries then falls back to the next. Works with any configured provider keys.
//
// OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./fallback
// GROQ_API_KEY=gsk-...  MISTRAL_API_KEY=...          go run ./fallback
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/providers"
	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func main() {
	configured := shared.ConfiguredProviders()
	if len(configured) == 0 {
		log.Fatal("No provider key set. Set at least one of: OPENAI_API_KEY, ANTHROPIC_API_KEY, GROQ_API_KEY, GEMINI_API_KEY, MISTRAL_API_KEY, TOGETHER_API_KEY, COHERE_API_KEY, DEEPSEEK_API_KEY")
	}

	// Build targets in the order providers were detected.
	targets := make([]aigateway.Target, len(configured))
	for i, p := range configured {
		targets[i] = aigateway.Target{VirtualKey: p.Name()}
	}
	// Retry up to 3 times on the primary provider before falling back.
	targets[0].Retry = &aigateway.RetryConfig{Attempts: 3}

	gw, err := aigateway.New(aigateway.Config{
		Strategy: aigateway.StrategyConfig{Mode: aigateway.ModeFallback},
		Targets:  targets,
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	for _, p := range configured {
		gw.RegisterProvider(p)
		fmt.Printf("Registered: %s\n", p.Name())
	}

	// Use the first model from the primary provider.
	model := configured[0].SupportedModels()[0]
	req := providers.Request{
		Model: model,
		Messages: []providers.Message{
			{Role: "user", Content: "Say hello in one sentence."},
		},
		MaxTokens: shared.IntPtr(50),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("\nRouting request (model=%s, strategy=fallback)...\n", model)
	resp, err := gw.Route(ctx, req)
	if err != nil {
		cancel()
		log.Printf("All providers failed: %v", err)
		os.Exit(1) //nolint:gocritic
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))
}
