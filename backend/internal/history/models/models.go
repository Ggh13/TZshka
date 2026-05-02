package historyModels

import (
	"time"

	"github.com/google/uuid"
)

type Correction struct {
	ID            uuid.UUID `json:"id"`
	SessionID     uuid.UUID `json:"sessionId"`
	InputContent  string    `json:"inputContent"`
	ResponseData  any       `json:"responseData"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
}

type CreateCorrectionRequest struct {
	SessionID    string `json:"sessionId" binding:"required"`
	InputContent string `json:"inputContent" binding:"required"`
	ResponseData any   `json:"responseData" binding:"required"`
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