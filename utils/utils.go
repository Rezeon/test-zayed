package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"ziyad-test/model"
)

// utilities

func generateTraceID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("TRC-%s", string(b))
}

func SendError(w http.ResponseWriter, code int, message string, errCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(model.CustomErrorResponse{
		Message:        message,
		ZiyadErrorCode: errCode,
		TraceID:        generateTraceID(),
	})
}
