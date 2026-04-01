package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/valuationservice/presentation/api"
)

func NewRouter(handler *presentation.ValuationController) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	setupPaths(r, handler)
	return r
}

func setupPaths(r *chi.Mux, handler *presentation.ValuationController) {
	r.Route("/valuation", func(router chi.Router) {
		router.Get("/", handler.GetValuationsHandler) //get all valuations
		router.Get("/{id}", handler.GetValuationHandler)
		//router.Patch("/{id}", handler.UpdateValuationHandler)  //update the calcualations
		router.Delete("/{id}", handler.DeleteValuationHandler) //update the calcualations
	})

}
