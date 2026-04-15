package application

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/claimservice/domain"
	"github.com/janicaleksander/cloud/claimservice/infrastructure/aws/utils"
	"github.com/janicaleksander/cloud/common/event"
)

type ClaimService struct {
	claimRepository    domain.ClaimRepository // //under the hood there is ref do persistance
	metaFileRepository domain.MetaFileRepository
	publisher          ClaimEventPublisher
	fileStorage        FileStorage
}

type ClaimEventPublisher interface {
	Publish(exchange string, msg interface{}) error
}
type FileStorage interface {
	StoreFile(ctx context.Context, bucket string, fileID string, contentType string, reader io.Reader) error
	GetFile(ctx context.Context, bucket string, fileID string) (io.ReadCloser, error)
}

func NewClaimService(claimRepo domain.ClaimRepository, metaRepo domain.MetaFileRepository, publisher ClaimEventPublisher, fileStorage FileStorage) *ClaimService {
	slog.Info("Creating ClaimService")
	return &ClaimService{
		claimRepository:    claimRepo,
		metaFileRepository: metaRepo,
		publisher:          publisher,
		fileStorage:        fileStorage,
	}
}

// TODO: user data - raz wykonywane przy budowie na root (nie trzeba sudo)

func DetectContentType(file *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return contentType, nil
}
func (c *ClaimService) CreateClaim(claim *domain.Claim, objectFiles []*os.File) (*domain.Claim, error) {
	slog.Info("Creating claim for user: ", "userID", claim.UserID, "vin", claim.VIN, "accidentDate", claim.AccidentDate)
	domainFiles := make([]*domain.File, 0, len(objectFiles))
	for idx := range objectFiles {
		fileUUID := uuid.New().String()
		newURL := utils.S3URL("us-east-1", "claim-cloud-bucket", fileUUID) + filepath.Ext(objectFiles[idx].Name())
		newKey := fileUUID + filepath.Ext(objectFiles[idx].Name())
		domainFile := &domain.File{
			StorageURL: newURL,
		}
		domainFiles = append(domainFiles, domainFile)
		contentType, err := DetectContentType(objectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
			//todo
		}
		err = c.fileStorage.StoreFile(context.Background(), "claim-cloud-bucket", newKey, contentType, objectFiles[idx])
		if err != nil {
			slog.Error(err.Error())
			//TODO sth
		}
		info, err := objectFiles[idx].Stat()
		if err != nil {
			slog.Error(err.Error())

			// handle error
			//TODO
		}

		size := info.Size()
		fmt.Println("size", size)
		metaFileModel := &domain.MetaFile{
			ID:       newKey,
			FileName: objectFiles[idx].Name(),
			FileExt:  filepath.Ext(objectFiles[idx].Name()),
			FileSize: float64(size),
			Date:     time.Now(),
			FileURL:  newURL,
		}

		_, err = c.metaFileRepository.Create(context.Background(), metaFileModel)
		if err != nil {
			slog.Error(err.Error())
			//todo
		}
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

func ParseS3URL(rawURL string) (bucket, key string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	// Host: "{bucket}.s3.{region}.amazonaws.com"
	host := u.Hostname()
	parts := strings.SplitN(host, ".", 2)
	if len(parts) < 2 || !strings.HasSuffix(parts[1], ".amazonaws.com") {
		return "", "", fmt.Errorf("not an S3 URL: %s", host)
	}

	bucket = parts[0]
	key = strings.TrimPrefix(u.Path, "/")

	return bucket, key, nil
}
func (c *ClaimService) GetFileFromStorage(fileID uint) (io.ReadCloser, *domain.MetaFile, error) {
	fileModel, err := c.claimRepository.GetFileById(context.Background(), fileID)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(fileModel.StorageURL)
	bucket, key, err := ParseS3URL(fileModel.StorageURL)
	fmt.Println(key)
	if err != nil {
		return nil, nil, err
	}
	metaFileDomain, err := c.metaFileRepository.GetFileById(context.Background(), key)
	if err != nil {
		return nil, nil, err
	}
	reader, err := c.fileStorage.GetFile(context.Background(), bucket, key)
	if err != nil {
		return nil, nil, err
	}
	return reader, metaFileDomain, nil
}
