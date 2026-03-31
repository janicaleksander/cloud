package presentation

import "time"

type CreateClaimRequestDTO struct {
	UserID       uint      `json:"user_id"`
	AccidentDate time.Time `json:"accident_date"`
	VIN          string    `json:"vin"`
}

type GetClaimResponseDTO struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	VIN          string            `json:"vin"`
	AccidentDate time.Time         `json:"accident_date"`
	Status       string            `json:"status"`
	Files        []FileResponseDTO `json:"files"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
type FileResponseDTO struct {
	ID         uint   `json:"id"`
	FileName   string `json:"file_name"`
	FileExt    string `json:"file_ext"`
	StorageURL string `json:"storage_url"`
}

type UpdateClaimRequestDTO struct {
	UserID uint `json:"user_id,omitempty"`
}
