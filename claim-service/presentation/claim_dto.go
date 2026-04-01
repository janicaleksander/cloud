package presentation

import "time"

type CreateClaimRequestDTO struct {
	UserID       uint      `json:"user_id"`
	Email        string    `json:"email"`
	AccidentDate time.Time `json:"accident_date"`
	VIN          string    `json:"vin"`
}

type GetClaimResponseDTO struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	Email        string            `json:"email"`
	VIN          string            `json:"vin"`
	AccidentDate time.Time         `json:"accident_date"`
	Status       string            `json:"status"`
	Files        []FileResponseDTO `json:"files"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
type FileResponseDTO struct {
	ID         uint   `json:"id"`
	FileName   string `json:"file_name"`
	FileExt    string `json:"file_ext"`
	StorageURL string `json:"storage_url"`
}

type UpdateClaimRequestDTO struct {
	Email string `json:"email,omitempty"`
}
