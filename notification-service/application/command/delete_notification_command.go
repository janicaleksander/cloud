package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteNotificationCommand struct {
	NotificationID string
}

type DeleteNotificationCommandHandler struct {
	repo domain.NotificationRepository
}

func NewDeleteNotificationCommandHandler(repo domain.NotificationRepository) *DeleteNotificationCommandHandler {
	return &DeleteNotificationCommandHandler{repo: repo}
}

func (h *DeleteNotificationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeleteNotificationCommand, *mediatr.Unit](h)
}

func (h *DeleteNotificationCommandHandler) Handle(ctx context.Context, cmd *DeleteNotificationCommand) (*mediatr.Unit, error) {
	nid, err := uuid.Parse(cmd.NotificationID)
	if err != nil {
		return nil, err
	}
	err = h.repo.DeleteNotificationByID(ctx, nid)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
