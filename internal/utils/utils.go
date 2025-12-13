package utils

import (
	"encoding/json"
	"net/http"
)

func JsonWriteWithBackup(w http.ResponseWriter, httpcode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpcode)
	json.NewEncoder(w).Encode(data)
}

func JsonError(w http.ResponseWriter, httpcode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpcode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": errorMsg,
	})
}
