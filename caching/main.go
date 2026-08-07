// Package main demonstrates the built-in response-cache plugin.
//
// response-cache stores completions keyed on the request, so an identical
// request is served from memory without calling the provider — saving latency
// and token cost. The plugin runs at two stages: before_request (serve a hit
// and skip the provider) and after_request (store the response). Both entries
// must carry identical config so they resolve to one shared cache.
//
// This example sends the same request twice and prints the latency of each: the
// second call is a cache hit — near-instant, and it returns the same response
// ID as the first, proving it came from the cache rather than the provider.
//
// OPENAI_API_KEY=sk-...        go run ./caching
// ANTHROPIC_API_KEY=sk-ant-... go run ./caching
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway-examples/shared"
	"github.com/ferro-labs/ai-gateway/config"
	"github.com/ferro-labs/ai-gateway/providers"

	// Register the built-in response-cache plugin factory via its init().
	_ "github.com/ferro-labs/ai-gateway/plugin/cache"
)

func main() {
	provider := shared.FirstProvider()
	model := shared.DefaultModel(provider)

	// The two stage entries MUST share identical config, or the gateway refuses
	// to start — config equality is what makes them one cache instance.
	cacheCfg := map[string]any{"max_age": 300, "max_entries": 1000}

	gw, err := aigateway.New(config.Config{
		Strategy: config.StrategyConfig{Mode: config.ModeSingle},
		Targets:  []config.Target{{VirtualKey: provider.Name()}},
		Plugins: []config.PluginConfig{
			{Name: "response-cache", Type: "transform", Stage: "before_request", Enabled: true, Config: cacheCfg},
			{Name: "response-cache", Type: "transform", Stage: "after_request", Enabled: true, Config: cacheCfg},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(provider)
	if err := gw.LoadPlugins(); err != nil {
		log.Fatalf("Failed to load plugins: %v", err)
	}

	req := providers.Request{
		Model:     model,
		Messages:  []providers.Message{{Role: "user", Content: "Name one primary color."}},
		MaxTokens: shared.IntPtr(10),
	}

	fmt.Printf("Provider: %s  Model: %s\n\n", provider.Name(), model)

	// First call: misses the cache, hits the provider, gets stored.
	id1, ms1 := routeTimed(gw, req)
	fmt.Printf("Call 1 (cache miss): %6.1f ms  id=%s\n", ms1, id1)

	// Second identical call: served from the cache, no provider call.
	id2, ms2 := routeTimed(gw, req)
	fmt.Printf("Call 2 (cache hit):  %6.1f ms  id=%s\n", ms2, id2)

	fmt.Println()
	if id1 == id2 {
		fmt.Printf("Same response ID both times → call 2 was served from cache (%.0fx faster).\n", ms1/max(ms2, 0.01))
	} else {
		fmt.Println("Response IDs differ — cache did not hit (check plugin config).")
	}
}

// routeTimed routes one request and returns its response ID and latency in ms.
func routeTimed(gw *aigateway.Gateway, req providers.Request) (id string, ms float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	start := time.Now()
	resp, err := gw.Route(ctx, req)
	ms = float64(time.Since(start).Microseconds()) / 1000.0
	if err != nil {
		log.Fatalf("Route failed: %v", err)
	}
	return resp.ID, ms
}
