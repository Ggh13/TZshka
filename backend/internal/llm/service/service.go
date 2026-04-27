package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

type LLMRequest struct {
	Mode     string `json:"mode"`
	Standard string `json:"standard,omitempty"`
	Content  string `json:"content"`
}

type LLMResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func (s *Service) ProcessText(ctx context.Context, mode, standard, content, sessionID string) (*LLMResponse, error) {
	s.log.Info("processing text", zap.String("mode", mode), zap.String("standard", standard), zap.String("session_id", sessionID))

	reqBody := LLMRequest{
		Mode:     mode,
		Standard: standard,
		Content:  content,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		s.log.Error("failed to marshal request", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.llmURL+"/llm/text", bytes.NewReader(body))
	if err != nil {
		s.log.Error("failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
		s.log.Error("failed to unmarshal response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	s.log.Info("text processed successfully")
	return &llmResp, nil
}
