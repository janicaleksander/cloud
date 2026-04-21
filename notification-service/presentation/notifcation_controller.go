package presentation

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/application/command"
	"github.com/janicaleksander/cloud/notificationservice/application/query"
	"github.com/mehdihadeli/go-mediatr"
)

type NotificationController struct {
}

func NewNotificationController() *NotificationController {
	slog.Info("Creating NotificationController")
	return &NotificationController{}
}
func success(w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if msg != nil {
		json.NewEncoder(w).Encode(msg)
	}
}

func successWithLocation(w http.ResponseWriter, location string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", location)
	w.WriteHeader(code)
}

func failure(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (nc *NotificationController) GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	slog.Info("HTTP GetNotificationsHandler called")
	q := GetNotificationsHTTPToQuery()
	notifications, err := mediatr.Send[*query.GetNotificationsQuery, *query.GetNotificationsQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notifications")
		return
	}
	success(w, map[string]any{"notifications": notifications.Notifications}, 200)

}
func (nc *NotificationController) GetNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	slog.Info("HTTP GetNotificationHandler called")
	idStr := chi.URLParam(r, "id")

	notificationID, err := uuid.Parse(idStr)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	q := GetNotificationHTTPToQuery(notificationID)
	notification, err := mediatr.Send[*query.GetNotificationQuery, *query.GetNotificationQueryResponse](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}
	success(w, map[string]any{"notification": notification}, 200)

}
func (nc *NotificationController) GetNotificationsForClaimIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	slog.Info("HTTP GetNotificationsForClaimIDHandler called")
	idStr := chi.URLParam(r, "id")
	claimID, err := uuid.Parse(idStr)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid claim ID")
		return
	}
	q := GetNotificationsForClaimIDHTTPToQuery(claimID)
	notificationForClaimID, err := mediatr.Send[*query.GetNotificationsForClaimIDQuery, *query.GetNotificationsForClaimIDQueryResult](context.Background(), q)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}
	success(w, map[string]any{"notifications": notificationForClaimID.Notifications}, 200)

}
func (nc *NotificationController) DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	slog.Info("HTTP DeleteNotificationHandler called")
	idStr := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(idStr)
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	cmd := DeleteNotificationHTTPToCommand(notificationID)
	_, err = mediatr.Send[*command.DeleteNotificationCommand, *mediatr.Unit](context.Background(), cmd)
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}
	success(w, nil, http.StatusNoContent)
}
