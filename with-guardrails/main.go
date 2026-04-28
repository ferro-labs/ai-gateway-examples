// Package main demonstrates using guardrail plugins in a standalone examples repo.
//
// Guardrails run before the request reaches the provider: the word-filter
// rejects blocked phrases, and max-token enforces token/message limits.
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
	"strings"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/plugin"
	"github.com/ferro-labs/ai-gateway/providers"
	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func init() {
	plugin.RegisterFactory("word-filter", func() plugin.Plugin { return &wordFilter{} })
	plugin.RegisterFactory("max-token", func() plugin.Plugin { return &maxTokenGuardrail{} })
}

func main() {
	provider := shared.FirstProvider()
	model := provider.SupportedModels()[0]

	gw, err := aigateway.New(aigateway.Config{
		Strategy: aigateway.StrategyConfig{Mode: aigateway.ModeSingle},
		Targets:  []aigateway.Target{{VirtualKey: provider.Name()}},
		Plugins: []aigateway.PluginConfig{
			{
				Name:    "word-filter",
				Type:    "guardrail",
				Stage:   "before_request",
				Enabled: true,
				Config: map[string]interface{}{
					"blocked_words":  []string{"password", "secret", "api_key"},
					"case_sensitive": false,
				},
			},
			{
				Name:    "max-token",
				Type:    "guardrail",
				Stage:   "before_request",
				Enabled: true,
				Config: map[string]interface{}{
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

type wordFilter struct {
	blockedWords  []string
	caseSensitive bool
}

func (w *wordFilter) Name() string { return "word-filter" }

func (w *wordFilter) Type() plugin.PluginType { return plugin.TypeGuardrail }

func (w *wordFilter) Init(config map[string]interface{}) error {
	if words, ok := config["blocked_words"]; ok {
		switch list := words.(type) {
		case []interface{}:
			for _, word := range list {
				if s, ok := word.(string); ok {
					w.blockedWords = append(w.blockedWords, s)
				}
			}
		case []string:
			w.blockedWords = append(w.blockedWords, list...)
		}
	}
	if cs, ok := config["case_sensitive"].(bool); ok {
		w.caseSensitive = cs
	}
	return nil
}

func (w *wordFilter) Execute(_ context.Context, pctx *plugin.Context) error {
	if pctx.Request == nil || len(w.blockedWords) == 0 {
		return nil
	}

	for _, msg := range pctx.Request.Messages {
		content := msg.Content
		if !w.caseSensitive {
			content = strings.ToLower(content)
		}
		for _, word := range w.blockedWords {
			check := word
			if !w.caseSensitive {
				check = strings.ToLower(check)
			}
			if strings.Contains(content, check) {
				pctx.Reject = true
				pctx.Reason = "blocked word detected: " + word
				return nil
			}
		}
	}

	return nil
}

type maxTokenGuardrail struct {
	maxTokens   int
	maxMessages int
}

func (m *maxTokenGuardrail) Name() string { return "max-token" }

func (m *maxTokenGuardrail) Type() plugin.PluginType { return plugin.TypeGuardrail }

func (m *maxTokenGuardrail) Init(config map[string]interface{}) error {
	m.maxTokens = 4096
	if v, ok := config["max_tokens"]; ok {
		switch val := v.(type) {
		case float64:
			m.maxTokens = int(val)
		case int:
			m.maxTokens = val
		}
	}

	m.maxMessages = 100
	if v, ok := config["max_messages"]; ok {
		switch val := v.(type) {
		case float64:
			m.maxMessages = int(val)
		case int:
			m.maxMessages = val
		}
	}

	return nil
}

func (m *maxTokenGuardrail) Execute(_ context.Context, pctx *plugin.Context) error {
	if pctx.Request == nil {
		return nil
	}

	if pctx.Request.MaxTokens != nil && *pctx.Request.MaxTokens > m.maxTokens {
		pctx.Reject = true
		pctx.Reason = fmt.Sprintf("max_tokens %d exceeds limit of %d", *pctx.Request.MaxTokens, m.maxTokens)
		return nil
	}

	if len(pctx.Request.Messages) > m.maxMessages {
		pctx.Reject = true
		pctx.Reason = fmt.Sprintf("message count %d exceeds limit of %d", len(pctx.Request.Messages), m.maxMessages)
		return nil
	}

	return nil
}
