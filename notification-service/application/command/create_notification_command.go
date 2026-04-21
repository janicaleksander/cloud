package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateNotificationCommand struct {
	ClaimID uuid.UUID
	Body    string
	SentTo  string
	Time    time.Time
}

type CreateNotificationCommandHandler struct {
	repo domain.NotificationRepository
}

func NewCreateNotificationCommandHandler(repo domain.NotificationRepository) *CreateNotificationCommandHandler {
	return &CreateNotificationCommandHandler{repo: repo}
}
func (h *CreateNotificationCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CreateNotificationCommand, *mediatr.Unit](h)
}

func (h *CreateNotificationCommandHandler) Handle(ctx context.Context, cmd *CreateNotificationCommand) (*mediatr.Unit, error) {
	notification := domain.NewNotification(
		uuid.New(),
		cmd.ClaimID,
		cmd.Body,
		cmd.SentTo,
		cmd.Time,
	)
	_, err := h.repo.SaveNotification(ctx, notification)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil
}
