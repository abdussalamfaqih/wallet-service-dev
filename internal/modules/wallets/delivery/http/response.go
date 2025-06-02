package http

import "net/http"

type (
	ResponsePayload struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func NewSuccessResponse(msg string, data interface{}) ResponsePayload {
	return ResponsePayload{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
}

func NewBadRequestResponse(msg string) ResponsePayload {
	return ResponsePayload{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NewInternalErrorResponse(msg string) ResponsePayload {
	return ResponsePayload{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}
