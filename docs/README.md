# ai-gemini — documentation

Google Gemini driver for togo ai

## Overview

Package gemini is a Google Gemini driver for togo ai (generateContent + embeddings).
Blank-import it and set AI_DRIVER=gemini + GEMINI_API_KEY.

## Install

```bash
togo install togo-framework/ai-gemini
```

Set `AI_DRIVER=gemini`.

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `GEMINI_API_KEY"` |
| `GOOGLE_API_KEY"` |

## Usage

```go
provider := ai.FromKernel(k)
resp, err := provider.Chat(ctx, []ai.Message{{Role: "user", Content: "Hello"}}, ai.Options{})
// streaming + provider.Embed(ctx, texts) for vectors; resp.Usage carries token counts
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-gemini
- Full README: ../README.md
