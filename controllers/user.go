package controllers

import (
	"api-go/models"
	characterModel "api-go/models/character"
	"api-go/models/entity"
	"api-go/services"
	userCharacterService "api-go/services/user_character"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Solicitação de cadastro de usuário")
	id, err := services.Create(user)
	var status = http.StatusOK
	var resp map[string]any
	if err != nil {
		status = http.StatusBadRequest
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar inserir: %v", err),
		}
	} else {
		resp = map[string]any{
			"error":   false,
			"message": fmt.Sprintf("Usuário criado com sucesso, ID: %d", id),
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
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	i := claims["id"].(float64)
	id := int(i)
	log.Printf("Solicitação de atualização de usuário")
	rows, err := services.Update(int64(id), user)
	var status = http.StatusOK
	var resp map[string]any
	if err != nil {
		status = http.StatusBadRequest
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar atualizar: %v", err),
		}
	} else {
		resp = map[string]any{
			"error":   false,
			"message": "Usuário atualizado com sucesso",
		}
	}
	if rows > 1 {
		log.Printf("Erro: foram atualizados %d registros", rows)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	i := claims["id"].(float64)
	id := int(i)
	log.Printf("Solicitação de detalhes de usuário")
	user, err := models.Get(int64(id))
	if err != nil {
		log.Printf("Erro ao buscar o usuário: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ucs, err := userCharacterService.List(int64(id))
	if err != nil {
		log.Printf("Erro ao obter personagens: %v", err)
	}
	var array []entity.UserCharacter
	if len(ucs) == 0 {
		array = []entity.UserCharacter{}
	} else {
		for _, uc := range ucs {
			c, err := characterModel.Get(uc.CharacterId)
			if err != nil {
				continue
			}
			uc.Character = c
			uc.HpMax = userCharacterService.CalculateHp(uc)
			array = append(array, uc)
		}
	}
	user.Characters = array
	w.Header().Add("Content-Type", "application/json")
	permission, err := models.GetPermission(user.ID)
	user.Permissions = permission
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
		"error":   false,
		"message": "Usuário removido com sucesso",
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Solicitação de autenticação de usuário")
	token, err := services.Auth(user)
	var status = http.StatusOK
	var resp map[string]any
	if err != nil {
		status = http.StatusForbidden
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar autenticar: %v", err),
		}
	} else {
		resp = map[string]any{
			"error":   false,
			"message": fmt.Sprintf("Usuário autenticado com sucesso"),
			"token":   token,
		}
		http.SetCookie(w, &http.Cookie{
			HttpOnly: true,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			SameSite: http.SameSiteLaxMode, // Uncomment below for HTTPS:
			// Secure: true,
			Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
			Value: token,
		})
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
