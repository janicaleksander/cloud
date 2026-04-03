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
	r.Use(middleware.RequestSize(10 << 20))

	setupPaths(r, claimHandler)
	return r
}
func setupPaths(r *chi.Mux, handler *presentation.ClaimController) {
	r.Route("/claim", func(router chi.Router) {
		router.Get("/", handler.GetClaimsHandler)          // all claims /claim
		router.Get("/{id}", handler.GetClaimHandler)       //claim /claim/id
		router.Post("/", handler.CreateClaimHandler)       //add claim
		router.Delete("/{id}", handler.DeleteClaimHandler) //delete claim /claim/id
		router.Patch("/{id}", handler.UpdateClaimHandler)  //delete claim /claim/id
	})
}
