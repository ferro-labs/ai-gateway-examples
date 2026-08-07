# Embeddings Example

Generates text embeddings through the gateway with `gw.Embed`. Embeddings turn text into vectors so you can measure semantic similarity — the core building block for search and RAG.

## Run

```bash
OPENAI_API_KEY=sk-... go run ./embeddings
```

Embeddings need an embedding-capable provider; this example uses OpenAI (`text-embedding-3-small`). Cohere, Mistral, and Gemini work identically — swap the provider and the model constant.

## What it does

1. Embeds three sentences in one `gw.Embed` call (two related, one unrelated)
2. Prints the vector dimensions
3. Prints the cosine similarity between the sentences, showing that related text scores much higher

## Expected output

```
Model: text-embedding-3-small   vector dimensions: 1536

  [0] The cat sat on the mat.
  [1] A feline rested on the rug.
  [2] Quarterly revenue grew 12%.

Cosine similarity:
  [0] vs [1] (related):   0.62
  [0] vs [2] (unrelated): 0.08
```

(Exact scores vary slightly by model version.) The same `gw.Embed` call powers a real retrieval pipeline: embed your documents once, embed the query at search time, and rank by cosine similarity.
