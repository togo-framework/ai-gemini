// Package gemini is a Google Gemini driver for togo ai (generateContent + embeddings).
// Blank-import it and set AI_DRIVER=gemini + GEMINI_API_KEY.
package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/ai"
	"github.com/togo-framework/togo"
)

const (
	base         = "https://generativelanguage.googleapis.com/v1beta"
	defaultModel = "gemini-1.5-flash"
	embedModel   = "text-embedding-004"
)

func init() {
	ai.RegisterDriver("gemini", func(k *togo.Kernel) (ai.Provider, error) {
		key := os.Getenv("GEMINI_API_KEY")
		if key == "" {
			key = os.Getenv("GOOGLE_API_KEY")
		}
		if key == "" {
			return nil, errors.New("ai-gemini: GEMINI_API_KEY not set")
		}
		return &provider{key: key, model: defaultModel, client: &http.Client{Timeout: 120 * time.Second}}, nil
	})
}

type provider struct {
	key, model string
	client     *http.Client
}

func (p *provider) post(ctx context.Context, path string, body, out any) error {
	buf, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%s?key=%s", base, path, p.key)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("ai-gemini: %s: %s", resp.Status, string(data))
	}
	return json.Unmarshal(data, out)
}

func (p *provider) Chat(ctx context.Context, req ai.ChatRequest) (ai.ChatResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}
	type part struct {
		Text string `json:"text"`
	}
	type content struct {
		Role  string `json:"role,omitempty"`
		Parts []part `json:"parts"`
	}
	var contents []content
	var sys strings.Builder
	for _, m := range req.Messages {
		if m.Role == "system" || m.Role == ai.RoleSystem {
			sys.WriteString(m.Content + "\n")
			continue
		}
		role := "user"
		if m.Role == ai.RoleAssistant || m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, content{Role: role, Parts: []part{{Text: m.Content}}})
	}
	body := map[string]any{"contents": contents}
	if strings.TrimSpace(sys.String()) != "" {
		body["systemInstruction"] = content{Parts: []part{{Text: strings.TrimSpace(sys.String())}}}
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []part `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}
	if err := p.post(ctx, "models/"+model+":generateContent", body, &out); err != nil {
		return ai.ChatResponse{}, err
	}
	var sb strings.Builder
	if len(out.Candidates) > 0 {
		for _, pt := range out.Candidates[0].Content.Parts {
			sb.WriteString(pt.Text)
		}
	}
	return ai.ChatResponse{
		Content: sb.String(),
		Model:   model,
		Usage:   ai.Usage{PromptTokens: out.UsageMetadata.PromptTokenCount, CompletionTokens: out.UsageMetadata.CandidatesTokenCount, TotalTokens: out.UsageMetadata.TotalTokenCount},
	}, nil
}

func (p *provider) Embed(ctx context.Context, req ai.EmbedRequest) (ai.EmbedResponse, error) {
	var res ai.EmbedResponse
	for _, in := range req.Inputs {
		body := map[string]any{"content": map[string]any{"parts": []map[string]string{{"text": in}}}}
		var out struct {
			Embedding struct {
				Values []float32 `json:"values"`
			} `json:"embedding"`
		}
		if err := p.post(ctx, "models/"+embedModel+":embedContent", body, &out); err != nil {
			return ai.EmbedResponse{}, err
		}
		res.Vectors = append(res.Vectors, out.Embedding.Values)
	}
	return res, nil
}
