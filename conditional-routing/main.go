// Package main demonstrates conditional routing: send different models to
// different providers through one gateway endpoint.
//
// This is the building block for cost/tier routing — point a cheap model at a
// cheap provider and a premium model at a premium one, then let callers pick by
// model name without knowing which provider serves it.
//
// The gateway matches each request's "model" field against the configured
// conditions and routes to the matching target; anything unmatched falls
// through to the default (the first target).
//
// Requires at least two provider keys, e.g.:
//
//	OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... go run ./conditional-routing
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/config"
	"github.com/ferro-labs/ai-gateway/providers"

	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func main() {
	configured := shared.ConfiguredProviders()
	if len(configured) < 2 {
		log.Fatal("Conditional routing needs at least 2 provider keys. Set any two of: OPENAI_API_KEY, ANTHROPIC_API_KEY, GROQ_API_KEY, GEMINI_API_KEY, MISTRAL_API_KEY, TOGETHER_API_KEY, COHERE_API_KEY, DEEPSEEK_API_KEY")
	}

	// Route each provider's own model to that provider. Both targets are
	// registered; conditions decide which one serves a given model.
	p0, p1 := configured[0], configured[1]
	m0, m1 := shared.DefaultModel(p0), shared.DefaultModel(p1)

	gw, err := aigateway.New(config.Config{
		Strategy: config.StrategyConfig{
			Mode: config.ModeConditional,
			Conditions: []config.Condition{
				{Key: config.ConditionKeyModel, Value: m0, TargetKey: p0.Name()},
				{Key: config.ConditionKeyModel, Value: m1, TargetKey: p1.Name()},
			},
		},
		// The first target is also the default for any unmatched model.
		Targets: []config.Target{{VirtualKey: p0.Name()}, {VirtualKey: p1.Name()}},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(p0)
	gw.RegisterProvider(p1)

	fmt.Printf("Routing rules:\n  model %q -> %s\n  model %q -> %s\n\n", m0, p0.Name(), m1, p1.Name())

	// Send each model and confirm it landed on the expected provider.
	for _, model := range []string{m0, m1} {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := gw.Route(ctx, providers.Request{
			Model:     model,
			Messages:  []providers.Message{{Role: "user", Content: "Reply with just 'ok'."}},
			MaxTokens: shared.IntPtr(10),
		})
		cancel()
		if err != nil {
			fmt.Printf("  model=%-40s ERROR %v\n", model, err)
			continue
		}
		fmt.Printf("  model=%-40s -> served by %q\n", model, resp.Provider)
	}
}
