package historyModels

import (
	"time"

	"github.com/google/uuid"
)

type Correction struct {
<<<<<<< HEAD
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"sessionId"`
	InputContent string                 `json:"inputContent"`
	ResponseData map[string]interface{} `json:"responseData"`
	CreatedAt    *time.Time             `json:"createdAt,omitempty"`
=======
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
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
}

type CorrectionsResponse struct {
	Corrections []Correction `json:"corrections"`
}
<<<<<<< HEAD
=======

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
