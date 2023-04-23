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
	err := json.NewEncoder(w).Encode(services.GetVersion())
	if err != nil {
		return
	}
}
