package application

/*package application

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/infrastructure/aws/utils"
	"github.com/janicaleksander/cloud/common/event"
)

type ClaimService struct {
	claimRepository domain.ClaimRepository
	publisher       ClaimEventPublisher
	fileStorage     FileStorage
}

func NewClaimService(claimRepo domain.ClaimRepository, publisher ClaimEventPublisher, fileStorage FileStorage) *ClaimService {
	slog.Info("Creating ClaimService")
	return &ClaimService{
		claimRepository: claimRepo,
		publisher:       publisher,
		fileStorage:     fileStorage,
	}
}

// TODO: user data - raz wykonywane przy budowie na root (nie trzeba sudo)

func (c *ClaimService) CreateClaim(claim *domain.Claim, objectFiles []*os.File) (*domain.Claim, error) {
	slog.Info("Creating claim for user: ", "userID", claim.UserID, "vin", claim.VIN, "accidentDate", claim.AccidentDate)
	domainFiles := make([]*domain.File, 0, len(objectFiles))
	for idx := range objectFiles {
		fileUUID := uuid.New().String()
		ext := filepath.Ext(objectFiles[idx].Name())
		newKey := fileUUID + ext
		newURL := utils.S3URL("us-east-1", "claim-cloud-bucket", newKey)

		contentType, err := DetectContentType(objectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
		}
		err = c.fileStorage.StoreFile(context.Background(), "claim-cloud-bucket", newKey, contentType, objectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
		}

		info, err := objectFiles[idx].Stat()
		if err != nil {
			slog.Error(err.Error())
		}

		domainFiles = append(domainFiles, &domain.File{
			FileName:   objectFiles[idx].Name(),
			FileExt:    ext,
			FileSize:   info.Size(),
			UploadedAt: time.Now(),
			StorageURL: newURL,
		})
	}

	claim.Status = domain.NEW
	if len(domainFiles) != 0 {
		claim.Files = domainFiles
	}

	urls := make([]string, 0, len(claim.Files))
	for idx := range claim.Files {
		urls = append(urls, claim.Files[idx].StorageURL)
	}
	savedClaim, err := c.claimRepository.Save(context.Background(), claim)
	if err != nil {
		return nil, err
	}
	err = c.publisher.Publish("events", event.RegisterUserForNotificationEvent{
		ClaimID: savedClaim.ID,
		Email:   savedClaim.Email,
	})
	if err != nil {
		return nil, err
	}

	err = c.pushClaimSubmittedEvent(&event.ClaimSubmittedEvent{
		UserID:       savedClaim.UserID,
		ClaimID:      savedClaim.ID,
		VIN:          savedClaim.VIN,
		AccidentDate: savedClaim.AccidentDate,
		StorageURL:   urls,
	})

	if err != nil {
		return nil, err
	}
	return savedClaim, nil
}
func (c *ClaimService) pushClaimSubmittedEvent(e *event.ClaimSubmittedEvent) error {
	return c.publisher.Publish("events", *e)
}

func (c *ClaimService) GetClaim(id uint) (*domain.Claim, error) {
	slog.Info("Getting claim with ID: ", "claimID", id)
	return c.claimRepository.GetById(context.Background(), id)
}
func (c *ClaimService) GetClaims() ([]*domain.Claim, error) {
	slog.Info("Getting all claims")
	return c.claimRepository.GetAll(context.Background())
}

func (c *ClaimService) DeleteClaim(id uint) error {
	slog.Info("Deleting claim with ID: ", "claimID", id)
	return c.claimRepository.DeleteById(context.Background(), id)
}

func (c *ClaimService) UpdateClaim(oldClaimDomain *domain.Claim, newUserEmail string) (*domain.Claim, error) {
	slog.Info("Updating claim with ID: ", "claimID", oldClaimDomain.ID, "newEmail", newUserEmail)
	if newUserEmail != oldClaimDomain.Email && newUserEmail != "" {
		oldClaimDomain.Email = newUserEmail
	}

	updatedClaim, err := c.claimRepository.Update(context.Background(), oldClaimDomain)
	if err != nil {
		return nil, err
	}
	err = c.publisher.Publish("events", event.ChangeEmailForNotification{
		ClaimID: oldClaimDomain.ID,
		Email:   newUserEmail,
	})
	if err != nil {
		return nil, err
	}
	return updatedClaim, nil
}

//rabbit events methods

func (c *ClaimService) ChangeClaimStatus(claimID uint, newStatus domain.Status) error {
	slog.Info("Changing claim status for ClaimID: ", "claimID", claimID, "newStatus", newStatus)
	return c.claimRepository.UpdateStatus(context.Background(), claimID, newStatus)
}

func (c *ClaimService) GetFileFromStorage(fileID uint) (io.ReadCloser, *domain.File, error) {
	file, err := c.claimRepository.GetFileById(context.Background(), fileID)
	if err != nil {
		return nil, nil, err
	}
	bucket, key, err := ParseS3URL(file.StorageURL)
	if err != nil {
		return nil, nil, err
	}
	reader, err := c.fileStorage.GetFile(context.Background(), bucket, key)
	if err != nil {
		return nil, nil, err
	}
	return reader, file, nil
}
*/
