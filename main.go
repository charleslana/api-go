package main

import (
	"api-go/configs"
	"api-go/controllers"
	"api-go/services"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

func main() {
	err := configs.Load()
	if err != nil {
		panic(err)
	}
	r := chi.NewRouter()
	r.Route("/public", func(r chi.Router) {
		r.Post("/user", controllers.Create)
		r.Post("/auth", controllers.Auth)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(services.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["id"])))
		})
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(jwtauth.Verifier(services.GetTokenAuth()))
		r.Use(services.AuthInterceptor)
		r.Put("/{id}", controllers.Update)
		r.Delete("/{id}", controllers.Delete)
		r.Get("/", controllers.List)
		//r.Get("/{id}", controllers.Get)
		r.Get("/details", controllers.Get)
	})
	log.Printf("Server started on %s", configs.GetServerPort())
	err = http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), r)
	if err != nil {
		return
	}
}
