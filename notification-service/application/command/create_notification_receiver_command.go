package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/notificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateNotificationReceiverCommand struct {
	ClaimID string
	Email   string
}

type CreateNotificationReceiverCommandHandler struct {
	repo domain.NotificationRepository
}

func NewCreateNotificationReceiverCommandHandler(repo domain.NotificationRepository) *CreateNotificationReceiverCommandHandler {
	return &CreateNotificationReceiverCommandHandler{
		repo: repo,
	}
}

func (h *CreateNotificationReceiverCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CreateNotificationReceiverCommand, *mediatr.Unit](h)
}

func (h *CreateNotificationReceiverCommandHandler) Handle(ctx context.Context, cmd *CreateNotificationReceiverCommand) (*mediatr.Unit, error) {
	cid, err := uuid.Parse(cmd.ClaimID)
	if err != nil {
		return nil, err
	}
	nr := &domain.NotificationReceiver{
		ID:      uuid.New(),
		ClaimID: cid,
		Email:   cmd.Email,
	}

	_, err = h.repo.SaveNotificationReceiver(ctx, nr)
	if err != nil {
		return nil, err
	}
	return &mediatr.Unit{}, nil

}
