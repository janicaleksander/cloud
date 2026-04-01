package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/decisionservice/presentation"
)

func NewRouter(handler *presentation.DecisionController) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	setupPaths(r, handler)
	return r
}

func setupPaths(r *chi.Mux, handler *presentation.DecisionController) {
	r.Route("/decision", func(router chi.Router) {
		router.Get("/", handler.GetDecisionsHandler)                //all  decisions
		router.Get("/{id}", handler.GetDecisionHandler)             //one decisions
		router.Get("/waiting", handler.GetWaitingDecisionsHandler)  //all waiting  decisions
		router.Post("/waiting/{id}", handler.UpdateDecisionHandler) // accept or deny the waiting decision
		router.Delete("/{id}", handler.DeleteDecisionHandler)       // del any decision
	})

}
