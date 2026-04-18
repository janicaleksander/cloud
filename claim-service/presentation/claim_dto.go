package presentation

import "time"

type CreateClaimRequestDTO struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	AccidentDate time.Time `json:"accident_date"`
	VIN          string    `json:"vin"`
}

type GetClaimResponseDTO struct {
	ID           string            `json:"id"`
	UserID       string            `json:"user_id"`
	Email        string            `json:"email"`
	VIN          string            `json:"vin"`
	AccidentDate time.Time         `json:"accident_date"`
	Status       string            `json:"status"`
	Files        []FileResponseDTO `json:"files"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
type FileResponseDTO struct {
	ID         string    `json:"id"`
	FileName   string    `json:"file_name"`
	FileExt    string    `json:"file_ext"`
	FileSize   int64     `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
	StorageURL string    `json:"storage_url"`
}

type UpdateClaimRequestDTO struct {
	Email string `json:"email,omitempty"`
}
