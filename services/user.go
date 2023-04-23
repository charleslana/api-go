package services

import (
	"api-go/models"
	"api-go/models/entity"
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
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
		err = fmt.Errorf("já existe um e-mail cadastrado")
		return 0, err
	}
	id, err = models.Insert(user)
	return id, err
}

func MakeToken(id int64, email string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"id": id, "email": email})
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
	token = MakeToken(user.ID, user.Email)
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
