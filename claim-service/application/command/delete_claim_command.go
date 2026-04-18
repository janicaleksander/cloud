package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/interfaces"
	"github.com/janicaleksander/cloud/claimservice/application/query"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteClaimCommand struct {
	ClaimID string
}

type DeleteClaimCommandHandler struct {
	repo        domain.ClaimRepository
	fileStorage interfaces.FileStorage
}

func NewDeleteClaimCommandHandler(r domain.ClaimRepository, f interfaces.FileStorage) *DeleteClaimCommandHandler {
	return &DeleteClaimCommandHandler{repo: r, fileStorage: f}
}

func (h *DeleteClaimCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*DeleteClaimCommand, *mediatr.Unit](h)
}
func (h *DeleteClaimCommandHandler) Handle(ctx context.Context, command *DeleteClaimCommand) (*mediatr.Unit, error) {
	cid, err := uuid.Parse(command.ClaimID)
	if err != nil {
		return nil, err
	}

	claim, err := h.repo.GetById(ctx, cid)
	if err != nil {
		return nil, err
	}

	for idx := range claim.Files {
		bucket, key, err := query.ParseS3URL(claim.Files[idx].StorageURL)
		if err != nil {
			return nil, err
		}
		err = h.fileStorage.RemoveFile(ctx, bucket, key)
		if err != nil {
			return nil, err
		}
	}

	err = h.repo.DeleteById(ctx, cid)
	if err != nil {
		return nil, err
	}

	return &mediatr.Unit{}, nil
}
