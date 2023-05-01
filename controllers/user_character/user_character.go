package controllers

import (
	characterModel "api-go/models/character"
	"api-go/models/entity"
	userCharacterService "api-go/services/user_character"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
	"strconv"
)

func checkSlot(ucs []entity.UserCharacter) bool {
	uc1 := ucs[0]
	uc2 := ucs[1]
	uc3 := ucs[2]
	if uc1.Slot != 1 || uc2.Slot != 2 || uc3.Slot != 3 {
		return true
	}
	return false
}

func duplicateInArray(arr []entity.UserCharacter) bool {
	visited := make(map[entity.UserCharacter]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] == true {
			return true
		} else {
			visited[arr[i]] = true
		}
	}
	return false
}

func Update(w http.ResponseWriter, r *http.Request) {
	var ucs []entity.UserCharacter
	err := json.NewDecoder(r.Body).Decode(&ucs)
	if err != nil {
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(ucs) != 3 {
		log.Printf("Quantia da lista inválido")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	check := checkSlot(ucs)
	if check {
		log.Printf("Slots da lista inválido")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	check = duplicateInArray(ucs)
	if check {
		log.Printf("Lista duplicada")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	i := claims["id"].(float64)
	userId := int64(i)
	rows := int64(0)
	log.Printf("Solicitação de atualização do personagem do usuário")
	if err != nil {
		log.Printf("Ocorreu um erro ao tentar atualizar o personagem do usuário: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var ids []int64
	for _, uc := range ucs {
		row, err := userCharacterService.Get(uc.ID, userId)
		if err != nil {
			continue
		} else {
			ids = append(ids, row.ID)
		}
		row.Slot = uc.Slot
		rows, err = userCharacterService.Update(userId, row)
		if err != nil {
			continue
		}
	}
	_, err = userCharacterService.ClearAllSlot(userId, ids)
	var status = http.StatusOK
	var resp map[string]any
	if err != nil {
		status = http.StatusBadRequest
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar atualizar o personagem do usuário: %v", err),
		}
	} else {
		resp = map[string]any{
			"error":   false,
			"message": "Personagem do usuário atualizado com sucesso",
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parse do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	i := claims["id"].(float64)
	userId := int64(i)
	log.Printf("Solicitação de detalhes do personagem do usuário")
	uc, err := userCharacterService.Get(int64(id), userId)
	if err != nil {
		var resp map[string]any
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar buscar o personagem do usuário: %v", err),
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
		return
	}
	c, err := characterModel.Get(uc.CharacterId)
	if err != nil {
		var resp map[string]any
		resp = map[string]any{
			"error":   true,
			"message": fmt.Sprintf("Ocorreu um erro ao tentar buscar o personagem: %v", err),
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			return
		}
		return
	}
	uc.Character = c
	w.Header().Add("Content-Type", "application/json")
	uc.HpMax = userCharacterService.CalculateHp(uc)
	err = json.NewEncoder(w).Encode(uc)
	if err != nil {
		return
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	i := claims["id"].(float64)
	userId := int64(i)
	log.Printf("Solicitação de todos os personagens do usuário")
	ucs, err := userCharacterService.List(userId)
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
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(array)
	if err != nil {
		return
	}
}
