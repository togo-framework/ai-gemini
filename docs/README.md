# ai-gemini — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Overview

Package gemini is a Google Gemini driver for togo ai (generateContent + embeddings).
Blank-import it and set AI_DRIVER=gemini + GEMINI_API_KEY.

## Install

```bash
togo install togo-framework/ai-gemini
```

Set `AI_DRIVER=gemini`.

## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `G` | _see provider docs_ |
| `GEMINI_API_KEY` | _see provider docs_ |
| `GOOGLE_API_KEY` | _see provider docs_ |

## Usage

```go
provider := ai.FromKernel(k)
resp, err := provider.Chat(ctx, []ai.Message{{Role: "user", Content: "Hello"}}, ai.Options{})
// streaming + provider.Embed(ctx, texts) for vectors; resp.Usage carries token counts
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-gemini
- README: ../README.md
