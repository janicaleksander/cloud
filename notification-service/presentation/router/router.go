package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/janicaleksander/cloud/notificationservice/presentation"
)

func NewRouter(handler *presentation.NotificationController) http.Handler {
	slog.Info("Setting up router")
	r := chi.NewMux()
	r.Use(middleware.Logger)
	setupPaths(r, handler)
	return r
}

func setupPaths(r *chi.Mux, handler *presentation.NotificationController) {
	r.Route("/notification", func(router chi.Router) {
		router.Get("/", handler.GetNotificationsHandler)                       // all notifications
		router.Get("/{id}", handler.GetNotificationHandler)                    //notifications for id
		router.Delete("/{id}", handler.DeleteNotificationHandler)              //delete by notID
		router.Get("/claimID/{id}", handler.GetNotificationsForClaimIDHandler) //notifications for claimID
	})

}
