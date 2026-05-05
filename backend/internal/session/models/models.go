package sessionModels

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatorID uuid.UUID  `json:"creatorId"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

type CreateSessionRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

type DeleteSessionRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

type JoinSessionRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

type SessionResponse struct {
	Session Session `json:"session"`
}

type SessionsResponse struct {
	Sessions []Session `json:"sessions"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
