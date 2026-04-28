// Package main demonstrates how to write and register a custom plugin with Ferro Labs AI Gateway.
//
// This example:
//   - Implements the plugin.Plugin interface with a request-ID injector plugin
//   - Registers it programmatically (no config file needed)
//   - Routes a request through the gateway so the plugin executes
//
// OPENAI_API_KEY=sk-...        go run ./custom-plugin
// ANTHROPIC_API_KEY=sk-ant-... go run ./custom-plugin
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/plugin"
	"github.com/ferro-labs/ai-gateway/providers"

	"github.com/ferro-labs/ai-gateway-examples/shared"
)

// requestIDPlugin stamps a trace ID onto every request via plugin.Context.Metadata.
type requestIDPlugin struct {
	prefix string
}

func (p *requestIDPlugin) Name() string            { return "request-id" }
func (p *requestIDPlugin) Type() plugin.PluginType { return plugin.TypeLogging }

func (p *requestIDPlugin) Init(config map[string]interface{}) error {
	if v, ok := config["prefix"].(string); ok && v != "" {
		p.prefix = v
	}
	return nil
}

func (p *requestIDPlugin) Execute(_ context.Context, pctx *plugin.Context) error {
	id := fmt.Sprintf("%s-%d", p.prefix, time.Now().UnixNano())
	pctx.Metadata["trace_id"] = id
	fmt.Printf("[request-id] trace_id=%s  model=%s\n", id, pctx.Request.Model)
	return nil
}

func main() {
	provider := shared.FirstProvider()
	model := provider.SupportedModels()[0]

	// 1. Create gateway with a single-provider strategy.
	gw, err := aigateway.New(aigateway.Config{
		Strategy: aigateway.StrategyConfig{Mode: aigateway.ModeSingle},
		Targets:  []aigateway.Target{{VirtualKey: provider.Name()}},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(provider)

	// 2. Instantiate and init the custom plugin.
	rp := &requestIDPlugin{}
	if err := rp.Init(map[string]interface{}{"prefix": "req"}); err != nil {
		log.Fatalf("Plugin init: %v", err)
	}

	// 3. Register it at the before_request stage.
	if err := gw.RegisterPlugin(plugin.StageBeforeRequest, rp); err != nil {
		log.Fatalf("Failed to register plugin: %v", err)
	}

	fmt.Printf("Provider: %s  Model: %s\n", provider.Name(), model)
	fmt.Println("Custom plugin registered. Routing request...")

	// 4. Route a request — the plugin fires automatically before the provider call.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := gw.Route(ctx, providers.Request{
		Model:     model,
		Messages:  []providers.Message{{Role: "user", Content: "Say 'plugin works!' and nothing else."}},
		MaxTokens: shared.IntPtr(10),
	})
	if err != nil {
		cancel()
		log.Printf("Route failed: %v", err)
		os.Exit(1) //nolint:gocritic
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))
}
