package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/policyverificationservice/presentation"
)

func NewRouter(policyHandler *presentation.PolicyController) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	setupPaths(r, policyHandler)
	return r
}

func setupPaths(r *chi.Mux, handler *presentation.PolicyController) {
	r.Route("/policy", func(router chi.Router) {
		router.Get("/", handler.GetPoliciesHandler)         //get all policies
		router.Get("/{id}", handler.GetPolicyHandler)       // get one policy
		router.Post("/", handler.CreatePolicyHandler)       //add policy
		router.Delete("/{id}", handler.DeletePolicyHandler) //delete policy
		router.Patch("/{id}", handler.UpdatePolicyHandler)  // update policy
	})

}
