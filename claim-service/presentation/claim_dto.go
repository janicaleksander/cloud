package presentation

import "time"

type CreateClaimRequestDTO struct {
	UserID uint            `json:"user_id"`
	CarID  uint            `json:"car_id"`
	Files  []CreateFileDTO `json:"files"`
}

type CreateFileDTO struct {
	FileName string `json:"file_name"`
	FileExt  string `json:"file_ext"`
}

type GetClaimResponseDTO struct {
	ID        uint              `json:"id"`
	UserID    uint              `json:"user_id"`
	CarID     uint              `json:"car_id"`
	Status    string            `json:"status"`
	Files     []FileResponseDTO `json:"files"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
type FileResponseDTO struct {
	ID       uint   `json:"id"`
	FileName string `json:"file_name"`
	FileExt  string `json:"file_ext"`
}
