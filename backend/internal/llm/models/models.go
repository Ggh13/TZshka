package llmModels

type TextRequest struct {
	Mode     string `json:"mode" form:"mode"`
	Standard string `json:"standard" form:"standard"`
	Content  string `json:"content" form:"content" binding:"required"`
}

type LLMSuccessResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    any  `json:"data"`
	Error   any  `json:"error"`
}

type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
