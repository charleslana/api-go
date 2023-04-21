package main

import (
	"api-go/configs"
	"api-go/controllers"
	"fmt"
	"github.com/go-chi/chi/v5"
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
	})
	r.Route("/user", func(r chi.Router) {
		r.Put("/{id}", controllers.Update)
		r.Delete("/{id}", controllers.Delete)
		r.Get("/", controllers.List)
		r.Get("/{id}", controllers.Get)
	})
	log.Printf("Server started on %s", configs.GetServerPort())
	err = http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), r)
	if err != nil {
		return
	}
}
