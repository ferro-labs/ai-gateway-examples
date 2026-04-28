// Package main demonstrates embedding Ferro Labs AI Gateway inside an existing Go HTTP server.
//
// Instead of running ferrogw as a standalone binary, you can import the library
// and mount the gateway handler alongside your own application routes.
//
// This example:
//   - Creates a standard net/http mux with existing application routes (/api/hello)
//   - Embeds gateway chat completions at /ai/v1/chat/completions
//   - Serves everything from a single server on :8080
//
// OPENAI_API_KEY=sk-...        go run ./embedded
// ANTHROPIC_API_KEY=sk-ant-... go run ./embedded
// # (any of the 8 supported provider keys work)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	aigateway "github.com/ferro-labs/ai-gateway"
	"github.com/ferro-labs/ai-gateway/providers"
	"github.com/ferro-labs/ai-gateway-examples/shared"
)

func main() {
	provider := shared.FirstProvider()

	// 1. Create and configure the gateway.
	gw, err := aigateway.New(aigateway.Config{
		Strategy: aigateway.StrategyConfig{Mode: aigateway.ModeSingle},
		Targets:  []aigateway.Target{{VirtualKey: provider.Name()}},
	})
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	gw.RegisterProvider(provider)

	// 2. Your existing application routes.
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "hello from your app"}) //nolint:errcheck,gosec
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK") //nolint:errcheck
	})

	// 3. Mount the gateway handler at a sub-path alongside your app routes.
	mux.HandleFunc("/ai/v1/chat/completions", gatewayHandler(gw))

	addr := ":8080"
	log.Printf("Provider: %s", provider.Name())
	log.Printf("Server listening on %s", addr)
	log.Println("  App route:     GET  /api/hello")
	log.Println("  Gateway route: POST /ai/v1/chat/completions")
	log.Fatal(http.ListenAndServe(addr, mux)) //nolint:gosec
}

// gatewayHandler returns an http.HandlerFunc that decodes a providers.Request,
// routes it through the gateway, and writes back a JSON response.
func gatewayHandler(gw *aigateway.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req providers.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		resp, err := gw.Route(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck,gosec
	}
}
