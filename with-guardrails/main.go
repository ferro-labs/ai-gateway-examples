// Package main demonstrates the gateway's built-in guardrail plugins.
//
// Guardrails run before a request reaches the provider. This example enables
// two that ship with the gateway, purely through configuration:
//
//   - word-filter: rejects requests whose content contains a blocked phrase
//   - max-token:   rejects requests that exceed token/message limits
//
// The blank imports below pull in the built-in plugin packages so their
// factories register at init time; the gateway then constructs them from the
// Plugins config and runs them via LoadPlugins. To write your own plugin
// instead, see the custom-plugin example.
//
// OPENAI_API_KEY=sk-...        go run ./with-guardrails
// ANTHROPIC_API_KEY=sk-ant-... go run ./with-guardrails
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway-examples/shared"
	"github.com/ferro-labs/ai-gateway/config"
	"github.com/ferro-labs/ai-gateway/providers"

	// Register the built-in guardrail plugin factories via their init().
	_ "github.com/ferro-labs/ai-gateway/plugin/maxtoken"
	_ "github.com/ferro-labs/ai-gateway/plugin/wordfilter"
)

func main() {
	provider := shared.FirstProvider()
	model := shared.DefaultModel(provider)

	gw, err := aigateway.New(config.Config{
		Strategy: config.StrategyConfig{Mode: config.ModeSingle},
		Targets:  []config.Target{{VirtualKey: provider.Name()}},
		Plugins: []config.PluginConfig{
			{
				Name:    "word-filter",
				Type:    "guardrail",
				Stage:   "before_request",
				Enabled: true,
				Config: map[string]any{
					"blocked_words":  []string{"password", "secret", "api_key"},
					"case_sensitive": false,
				},
			},
			{
				Name:    "max-token",
				Type:    "guardrail",
				Stage:   "before_request",
				Enabled: true,
				Config: map[string]any{
					"max_tokens":   4096,
					"max_messages": 50,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(provider)
	if err := gw.LoadPlugins(); err != nil {
		log.Fatalf("Failed to load plugins: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Provider: %s  Model: %s\n", provider.Name(), model)

	// This request passes through — no blocked words.
	fmt.Println("\n--- Request 1: clean message (should pass) ---")
	resp, err := gw.Route(ctx, providers.Request{
		Model:    model,
		Messages: []providers.Message{{Role: "user", Content: "Tell me a joke about Go programming."}},
	})
	if err != nil {
		fmt.Printf("Rejected: %v\n", err)
	} else {
		out, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(out))
	}

	// This request is blocked by the word-filter.
	fmt.Println("\n--- Request 2: blocked word 'password' (should be rejected) ---")
	_, err = gw.Route(ctx, providers.Request{
		Model:    model,
		Messages: []providers.Message{{Role: "user", Content: "What is a secure password strategy?"}},
	})
	if err != nil {
		fmt.Printf("Rejected: %v\n", err)
	} else {
		fmt.Println("Passed (word-filter may not block this phrasing)")
	}
}
