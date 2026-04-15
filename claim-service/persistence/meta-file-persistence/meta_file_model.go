package meta_file_persistence

import (
	"time"
)

type MetaFileModel struct {
	ID       string    `dynamodbav:"id"`
	FileName string    `dynamodbav:"fileName"`
	FileExt  string    `dynamodbav:"fileExt"`
	FileSize float64   `dynamodbav:"fileSize"`
	Date     time.Time `dynamodbav:"date"`
	FileURL  string    `dynamodbav:"fileURL"`
}
