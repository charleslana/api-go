package controllers

import (
	"api-go/services"
	"encoding/json"
	"log"
	"net/http"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Printf("Solicitação da versão do app")
	w.Header().Add("Content-Type", "application/json")
	resp := map[string]any{
		"error":   false,
		"message": services.GetVersion(),
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
