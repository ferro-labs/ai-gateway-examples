// Package main demonstrates streaming chat completions through the AI Gateway.
//
// The gateway returns a channel of StreamChunk objects. Each chunk contains
// incremental content (delta) that you print as it arrives — giving the user
// real-time output instead of waiting for the full response.
//
// OPENAI_API_KEY=sk-...        go run ./streaming
// ANTHROPIC_API_KEY=sk-ant-... go run ./streaming
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/providers"

	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func main() {
	provider := shared.FirstProvider()
	model := provider.SupportedModels()[0]

	gw, err := aigateway.New(aigateway.Config{
		Strategy: aigateway.StrategyConfig{Mode: aigateway.ModeSingle},
		Targets:  []aigateway.Target{{VirtualKey: provider.Name()}},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(provider)

	req := providers.Request{
		Model: model,
		Messages: []providers.Message{
			{Role: "user", Content: "Write a haiku about distributed systems."},
		},
		Stream:    true,
		MaxTokens: shared.IntPtr(100),
	}

	fmt.Printf("Provider: %s  Model: %s\n", provider.Name(), model)
	fmt.Println("Streaming response:")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ch, err := gw.RouteStream(ctx, req)
	if err != nil {
		cancel()
		log.Printf("Stream failed: %v", err)
		os.Exit(1) //nolint:gocritic
	}

	var totalTokens int
	for chunk := range ch {
		if chunk.Error != nil {
			fmt.Printf("\nStream error: %v\n", chunk.Error)
			break
		}
		for _, choice := range chunk.Choices {
			fmt.Print(choice.Delta.Content)
		}
		if chunk.Usage != nil {
			totalTokens = chunk.Usage.TotalTokens
		}
	}

	fmt.Println()
	if totalTokens > 0 {
		fmt.Printf("\nTotal tokens: %d\n", totalTokens)
	}
}
