package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/notificationservice/presentation"
)

func NewRouter(handler *presentation.NotificationController) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	setupPaths(r, handler)
	return r
}

func setupPaths(r *chi.Mux, handler *presentation.NotificationController) {
	r.Route("/notification", func(router chi.Router) {
		router.Get("/", handler.GetNotifications)                       // all notifications
		router.Get("/{id}", handler.GetNotification)                    //notifications for id
		router.Delete("/{id}", handler.DeleteNotification)              //delete by notID
		router.Get("/claimID/{id}", handler.GetNotificationsForClaimID) //notifications for claimID
	})

}
