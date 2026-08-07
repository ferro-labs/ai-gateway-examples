// Package main demonstrates generating embeddings through the gateway.
//
// Embeddings turn text into vectors so you can measure semantic similarity —
// the core building block for search and RAG. The gateway exposes them via
// gw.Embed, mirroring the OpenAI /v1/embeddings schema.
//
// To make the output meaningful, this example embeds three sentences and prints
// the cosine similarity between them: two related sentences score high, the
// unrelated one scores low.
//
// Embeddings need an embedding-capable provider; this example uses OpenAI.
// Cohere, Mistral, and Gemini work the same way — swap the provider and model.
//
// OPENAI_API_KEY=sk-... go run ./embeddings
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/config"
	"github.com/ferro-labs/ai-gateway/providers"
	openaipkg "github.com/ferro-labs/ai-gateway/providers/openai"
)

const embeddingModel = "text-embedding-3-small"

func main() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Fatal("embeddings example requires OPENAI_API_KEY (an embedding-capable provider)")
	}

	gw, err := aigateway.New(config.Config{
		Strategy: config.StrategyConfig{Mode: config.ModeSingle},
		Targets:  []config.Target{{VirtualKey: "openai"}},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	provider, err := openaipkg.New(key, "")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}
	gw.RegisterProvider(provider)

	sentences := []string{
		"The cat sat on the mat.",
		"A feline rested on the rug.", // related to #0
		"Quarterly revenue grew 12%.", // unrelated
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := gw.Embed(ctx, providers.EmbeddingRequest{
		Model: embeddingModel,
		Input: sentences,
	})
	if err != nil {
		log.Fatalf("Embed failed: %v", err)
	}
	if len(resp.Data) != len(sentences) {
		log.Fatalf("expected %d embeddings, got %d", len(sentences), len(resp.Data))
	}

	fmt.Printf("Model: %s   vector dimensions: %d\n\n", resp.Model, len(resp.Data[0].Embedding))
	for i, s := range sentences {
		fmt.Printf("  [%d] %s\n", i, s)
	}
	fmt.Printf("\nCosine similarity:\n")
	fmt.Printf("  [0] vs [1] (related):   %.4f\n", cosine(resp.Data[0].Embedding, resp.Data[1].Embedding))
	fmt.Printf("  [0] vs [2] (unrelated): %.4f\n", cosine(resp.Data[0].Embedding, resp.Data[2].Embedding))
}

// cosine returns the cosine similarity of two equal-length vectors, in [-1, 1].
// Higher means more semantically similar.
func cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
