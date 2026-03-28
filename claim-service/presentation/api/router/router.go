package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/claimservice/presentation"
)

func NewRouter(claimHandler *presentation.ClaimController) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	setupPaths(r, claimHandler)
	return r
}
func setupPaths(r *chi.Mux, handler *presentation.ClaimController) {
	r.Route("/claim", func(router chi.Router) {
		router.Get("/", nil)                         // all claims /claim
		router.Get("/{id}", handler.GetClaimHandler) //claim /claim/id
		router.Post("/", handler.CreateClaimHandler) //add claim
		router.Delete("/{id}", nil)                  //delete claim /claim/id
		router.Put("/{id}", nil)                     //delete claim /claim/id
	})
}
