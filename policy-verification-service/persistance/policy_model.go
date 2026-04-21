package persistance

import (
	"time"

	"github.com/google/uuid"
)

type PolicyModel struct {
	ID     uuid.UUID `dynamodbav:"id"`
	UserID uuid.UUID `dynamodbav:"user_id"`
	VIN    string    `dynamodbav:"vin"`
	From   time.Time `dynamodbav:"from"`
	To     time.Time `dynamodbav:"to"`
}
