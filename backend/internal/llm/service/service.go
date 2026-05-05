package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type Service struct {
	httpClient *http.Client
	llmURL     string
	log        *zap.Logger
}

func New(httpClient *http.Client, llmURL string, log *zap.Logger) *Service {
	return &Service{
		httpClient: httpClient,
		llmURL:     llmURL,
		log:        log,
	}
}

type LLMResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func (s *Service) ProcessText(ctx context.Context, mode, standard, content string) (*LLMResponse, error) {
	s.log.Info("processing text", zap.String("mode", mode), zap.String("standard", standard))

	data := url.Values{}
	data.Set("mode", mode)
	if standard != "" {
		data.Set("standard", standard)
	}
	data.Set("content", content)

	req, err := http.NewRequestWithContext(ctx, "POST", s.llmURL+"/llm/text", bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		s.log.Error("failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.log.Error("failed to call llm api", zap.Error(err))
		return nil, fmt.Errorf("failed to call llm api: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("failed to read response", zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(respBody, &llmResp); err != nil {
		s.log.Error("failed to unmarshal response", zap.Error(err), zap.String("body", string(respBody)))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.log.Info("text processed successfully")
	return &llmResp, nil
}
