package services

import (
	"api-go/models"
	characterModel "api-go/models/character"
	"api-go/models/entity"
	userCharacterService "api-go/services/user_character"
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"strings"
)

func Create(user entity.User) (id int64, err error) {
	if user.Email == "" {
		err = fmt.Errorf("email em branco")
		return 0, err
	}
	if user.Password == "" {
		err = fmt.Errorf("senha em branco")
		return 0, err
	}
	var rows int64
	rows, err = models.CountByEmail(user.Email)
	if rows > 0 {
		err = fmt.Errorf("já existe o e-mail cadastrado")
		return 0, err
	}
	id, err = models.Insert(user)
	for i := 1; i <= 4; i++ {
		c, err := characterModel.Get(int64(i))
		if err != nil {
			continue
		}
		slot := int8(i)
		if i == 4 {
			slot = 0
		}
		var uc = entity.UserCharacter{Level: 1, HpMin: c.Hp, CharacterId: int64(i), UserId: id, Slot: slot}
		_, err = userCharacterService.Create(uc)
		if err != nil {
			continue
		}
	}
	return id, err
}

func Update(id int64, user entity.User) (rows int64, err error) {
	n := strings.TrimSpace(*user.Name)
	if n == "" {
		err = fmt.Errorf("nome em branco")
		return 0, err
	}
	rows, err = models.CountByName(n)
	if rows > 0 {
		userExist, _ := models.GetByName(n)
		if userExist.ID != id && strings.EqualFold(n, *userExist.Name) {
			err = fmt.Errorf("já existe o nome cadastrado")
			return 0, err
		}
	}
	rows, err = models.Update(id, n)
	return rows, err
}

func MakeToken(user entity.User) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"id": user.ID, "email": user.Email, "permissions": user.Permissions})
	return tokenString
}

func Auth(user entity.User) (token string, err error) {
	if user.Email == "" || user.Password == "" {
		err = fmt.Errorf("email ou senha em branco")
		return "", err
	}
	user, err = models.GetUserDetails(user)
	if err != nil {
		err = fmt.Errorf("email ou senha inválidos")
		return "", err
	}
	permission, err := models.GetPermission(user.ID)
	user.Permissions = permission
	token = MakeToken(user)
	return token, nil
}

func AuthInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			resp := map[string]any{
				"error":   true,
				"message": err.Error(),
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err = json.NewEncoder(w).Encode(resp)
			//http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if token == nil || jwt.Validate(token) != nil {
			resp := map[string]any{
				"error":   true,
				"message": http.StatusText(http.StatusUnauthorized),
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			err = json.NewEncoder(w).Encode(resp)
			//http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AllowRoles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		//claims["permissions"] = []string{"user", "admin"}
		//fmt.Printf("%v", claims)
		if claims["permissions"] == nil {
			resp := map[string]any{
				"error":   true,
				"message": "Não permitido",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				return
			}
			return
		}
		permissions := claims["permissions"].([]interface{})
		if !containsRoles(permissions, "admin") {
			resp := map[string]any{
				"error":   true,
				"message": "Não permitido",
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				return
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func containsRoles(p []interface{}, str string) bool {
	for _, data := range p {
		for _, v := range data.(map[string]interface{}) {
			if v == str {
				return true
			}
		}
	}
	return false
}
