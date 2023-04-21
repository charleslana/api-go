package controllers

import (
	"api-go/models"
	"api-go/models/entity"
	"api-go/services"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := services.Create(user)
	var status = http.StatusOK
	var resp map[string]any
	if err != nil {
		status = http.StatusBadRequest
		resp = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Ocorreu um erro ao tentar inserir: %v", err),
		}
	} else {
		resp = map[string]any{
			"Error":   false,
			"Message": fmt.Sprintf("Usuário criado com sucesso, ID: %d", id),
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parse do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var user entity.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := models.Update(int64(id), *user.Name)
	if err != nil {
		log.Printf("Erro ao atualizar o usuário: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("Erro: foram atualizados %d registros", rows)
	}
	resp := map[string]any{
		"Error":   false,
		"Message": "Usuário atualizado com sucesso",
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parse do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user, err := models.Get(int64(id))
	if err != nil {
		log.Printf("Erro ao buscar o usuário: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAll()
	if err != nil {
		log.Printf("Erro ao obter usuários: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parse do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rows, err := models.Delete(int64(id))
	if err != nil {
		log.Printf("Erro ao remover o usuário: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if rows > 1 {
		log.Printf("Erro: foram removidos %d registros", rows)
	}
	resp := map[string]any{
		"Error":   false,
		"Message": "Usuário removido com sucesso",
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
