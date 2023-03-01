package controllers

import "encoding/json"

type ErrorResponse struct {
	Message string
}

func createErrorResponse(message string) []byte {
	resp := &ErrorResponse{
		Message: message,
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		return []byte{}
	}
	return bytes
}
