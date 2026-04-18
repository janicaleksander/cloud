package command

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/application/interfaces"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/infrastructure/aws/utils"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateClaimCommand struct {
	ID           string
	UserID       string
	Email        string
	VIN          string
	AccidentDate time.Time
	ObjectFiles  []*os.File
}

type CreateClaimCommandHandler struct {
	repo        domain.ClaimRepository
	fileStorage interfaces.FileStorage
	publisher   interfaces.ClaimEventPublisher
}

func NewCreateClaimCommandHandler(r domain.ClaimRepository, p interfaces.ClaimEventPublisher, f interfaces.FileStorage) *CreateClaimCommandHandler {
	return &CreateClaimCommandHandler{repo: r, publisher: p, fileStorage: f}
}
func (h *CreateClaimCommandHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CreateClaimCommand, *mediatr.Unit](h)
}

func (h *CreateClaimCommandHandler) Handle(ctx context.Context, command *CreateClaimCommand) (*mediatr.Unit, error) {

	domainFiles := make([]*domain.File, 0, len(command.ObjectFiles))
	for idx := range command.ObjectFiles {
		fileUUID := uuid.New()
		ext := filepath.Ext(command.ObjectFiles[idx].Name())
		newKey := fileUUID.String() + ext
		newURL := utils.S3URL("us-east-1", "claim-cloud-bucket", newKey)

		contentType, err := utils.DetectContentType(command.ObjectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
		}
		err = h.fileStorage.StoreFile(context.Background(), "claim-cloud-bucket", newKey, contentType, command.ObjectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
		}

		info, err := command.ObjectFiles[idx].Stat()
		if err != nil {
			slog.Error(err.Error())
		}

		domainFiles = append(domainFiles, &domain.File{
			ID:         fileUUID,
			FileName:   command.ObjectFiles[idx].Name(),
			FileExt:    ext,
			FileSize:   info.Size(),
			UploadedAt: time.Now(),
			StorageURL: newURL,
		})
	}
	claimDomain := CreateClaimCommandToDomain(command)
	claimDomain.Status = domain.NEW
	if len(domainFiles) != 0 {
		claimDomain.Files = domainFiles
	}

	urls := make([]string, 0, len(claimDomain.Files))
	for idx := range claimDomain.Files {
		urls = append(urls, claimDomain.Files[idx].StorageURL)
	}
	savedClaim, err := h.repo.Save(context.Background(), claimDomain)
	if err != nil {
		return nil, err
	}
	err = h.publisher.Publish("events", event.RegisterUserForNotificationEvent{
		ClaimID: savedClaim.ID.String(),
		Email:   savedClaim.Email,
	})
	if err != nil {
		return nil, err
	}

	err = h.publisher.Publish("events", &event.ClaimSubmittedEvent{
		UserID:       savedClaim.UserID.String(),
		ClaimID:      savedClaim.ID.String(),
		VIN:          savedClaim.VIN,
		AccidentDate: savedClaim.AccidentDate,
		StorageURL:   urls,
	})

	return &mediatr.Unit{}, err
}
