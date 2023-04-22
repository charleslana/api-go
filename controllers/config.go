package controllers

import (
	"api-go/services"
	"encoding/json"
	"net/http"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(services.GetVersion())
	if err != nil {
		return
	}
}
