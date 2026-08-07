// Package main demonstrates loading gateway configuration from a YAML file.
//
// Every other example builds config.Config in Go. Real deployments keep it in a
// file and load it with config.LoadConfig — edit YAML, restart, no recompile.
// This example loads gateway.yaml (strategy + targets + a "fast" alias), then
// routes a request using the alias instead of a hardcoded model ID.
//
// The target in gateway.yaml is "openai", so this example needs OPENAI_API_KEY.
// To use another provider, change the target's virtual_key and register that
// provider instead.
//
// OPENAI_API_KEY=sk-... go run ./config-file
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/config"
	"github.com/ferro-labs/ai-gateway/providers"
	openaipkg "github.com/ferro-labs/ai-gateway/providers/openai"
)

func main() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("config-file example requires OPENAI_API_KEY (gateway.yaml targets the openai provider)")
	}

	// Load the config from YAML sitting next to this source file.
	_, thisFile, _, _ := runtime.Caller(0)
	cfgPath := filepath.Join(filepath.Dir(thisFile), "gateway.yaml")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", cfgPath, err)
	}
	fmt.Printf("Loaded config: strategy=%s, targets=%d, aliases=%v\n",
		cfg.Strategy.Mode, len(cfg.Targets), cfg.Aliases)

	gw, err := aigateway.New(*cfg)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	provider, err := openaipkg.New(key, "")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}
	gw.RegisterProvider(provider)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ask for the "fast" alias — the gateway resolves it to gpt-4o-mini per the
	// aliases block, so the caller never hardcodes a provider's model string.
	fmt.Println("\nRouting request with model alias \"fast\"...")
	resp, err := gw.Route(ctx, providers.Request{
		Model:     "fast",
		Messages:  []providers.Message{{Role: "user", Content: "Say hello in one sentence."}},
		MaxTokens: intPtr(30),
	})
	if err != nil {
		log.Fatalf("Route failed: %v", err)
	}

	fmt.Printf("Alias \"fast\" resolved to model: %s\n", resp.Model)
	out, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(out))
}

func intPtr(i int) *int { return &i }
