package shared

import (
	"log"
	"os"

	"github.com/ferro-labs/ai-gateway/providers"

	anthropicpkg "github.com/ferro-labs/ai-gateway/providers/anthropic"
	coherepkg "github.com/ferro-labs/ai-gateway/providers/cohere"
	deepseekpkg "github.com/ferro-labs/ai-gateway/providers/deepseek"
	geminipkg "github.com/ferro-labs/ai-gateway/providers/gemini"
	groqpkg "github.com/ferro-labs/ai-gateway/providers/groq"
	mistralpkg "github.com/ferro-labs/ai-gateway/providers/mistral"
	openaipkg "github.com/ferro-labs/ai-gateway/providers/openai"
	togetherpkg "github.com/ferro-labs/ai-gateway/providers/together"
)

type providerEntry struct {
	env    string
	create func(key string) (providers.Provider, error)
}

var candidates = []providerEntry{
	{"OPENAI_API_KEY", func(k string) (providers.Provider, error) { return openaipkg.New(k, "") }},
	{"ANTHROPIC_API_KEY", func(k string) (providers.Provider, error) { return anthropicpkg.New(k, "") }},
	{"GROQ_API_KEY", func(k string) (providers.Provider, error) { return groqpkg.New(k, "") }},
	{"GEMINI_API_KEY", func(k string) (providers.Provider, error) { return geminipkg.New(k, "") }},
	{"MISTRAL_API_KEY", func(k string) (providers.Provider, error) { return mistralpkg.New(k, "") }},
	{"TOGETHER_API_KEY", func(k string) (providers.Provider, error) { return togetherpkg.New(k, "") }},
	{"COHERE_API_KEY", func(k string) (providers.Provider, error) { return coherepkg.New(k, "") }},
	{"DEEPSEEK_API_KEY", func(k string) (providers.Provider, error) { return deepseekpkg.New(k, "") }},
}

const noKeyMsg = "No provider key set. Set at least one of: OPENAI_API_KEY, ANTHROPIC_API_KEY, GROQ_API_KEY, GEMINI_API_KEY, MISTRAL_API_KEY, TOGETHER_API_KEY, COHERE_API_KEY, DEEPSEEK_API_KEY"

// FirstProvider returns the first provider for which an API key is set.
// It calls log.Fatal if no key is configured.
func FirstProvider() providers.Provider {
	for _, c := range candidates {
		if key := os.Getenv(c.env); key != "" {
			p, err := c.create(key)
			if err != nil {
				log.Fatalf("Failed to create provider for %s: %v", c.env, err)
			}
			return p
		}
	}
	log.Fatal(noKeyMsg)
	return nil
}

// ConfiguredProviders returns a provider instance for every key that is set.
func ConfiguredProviders() []providers.Provider {
	var result []providers.Provider
	for _, c := range candidates {
		if key := os.Getenv(c.env); key != "" {
			p, err := c.create(key)
			if err != nil {
				log.Fatalf("Failed to create provider for %s: %v", c.env, err)
			}
			result = append(result, p)
		}
	}
	return result
}

// IntPtr is a convenience helper for pointer-to-int fields.
func IntPtr(i int) *int { return &i }
