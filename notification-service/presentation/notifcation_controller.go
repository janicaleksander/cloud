package presentation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/janicaleksander/cloud/notificationservice/application"
)

type NotificationController struct {
	notificationService *application.NotificationService
}

func NewNotificationController(notificationService *application.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}
func success(w http.ResponseWriter, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(msg)
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
	notificationsDomain, err := nc.notificationService.GetNotifications()
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notifications")
		return
	}
	notificationsDTO := make([]GetNotificationResponseDTO, 0, len(notificationsDomain))
	for _, notification := range notificationsDomain {
		notificationsDTO = append(notificationsDTO, *GetNotificationDomainToResponse(notification))
	}
	success(w, map[string]any{"notifications": notificationsDTO})

}
func (nc *NotificationController) GetNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	idStr, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	notificationDomain, err := nc.notificationService.GetNotification(uint(idStr))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}
	notificationDTO := GetNotificationDomainToResponse(notificationDomain)
	success(w, map[string]any{"notification": notificationDTO})

}
func (nc *NotificationController) GetNotificationsForClaimIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	idStr, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	notificationDomain, err := nc.notificationService.GetNotificationsForClaimID(uint(idStr))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}
	notificationDTO := make([]GetNotificationResponseDTO, 0, len(notificationDomain))
	for _, notification := range notificationDomain {
		notificationDTO = append(notificationDTO, *GetNotificationDomainToResponse(notification))
	}
	success(w, map[string]any{"notifications": notificationDTO})

}
func (nc *NotificationController) DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		failure(w, http.StatusBadRequest, "Invalid method")
		return
	}
	idStr, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		failure(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}
	err = nc.notificationService.DeleteNotification(uint(idStr))
	if err != nil {
		failure(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}
	success(w, map[string]string{"message": "Notification deleted successfully"})
}
