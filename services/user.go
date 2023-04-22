package services

import (
	"api-go/models"
	"api-go/models/entity"
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
	id, err = models.Insert(user)
	return id, err
}

func MakeToken(email string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"email": email})
	return tokenString
}

func Auth(user entity.User) (token string, err error) {
	if user.Email == "" || user.Password == "" {
		err = fmt.Errorf("email ou senha em branco")
		return "", err
	}
	token = MakeToken(user.Email)
	return token, nil
}

func AuthInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())
		if token != nil && jwt.Validate(token) == nil {
			http.Redirect(w, r, "/profile", 302)
		}
		next.ServeHTTP(w, r)
	})
}
