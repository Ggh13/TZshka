package historyModels

import (
	"time"

	"github.com/google/uuid"
)

type Correction struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"sessionId"`
	InputContent string                 `json:"inputContent"`
	ResponseData map[string]interface{} `json:"responseData"`
	CreatedAt    *time.Time             `json:"createdAt,omitempty"`
}

type CorrectionsResponse struct {
	Corrections []Correction `json:"corrections"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
