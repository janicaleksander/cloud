package persistance

import (
	"time"
)

type PolicyModel struct {
	ID     string    `dynamodbav:"policy_id"`
	UserID string    `dynamodbav:"user_id"`
	VIN    string    `dynamodbav:"vin"`
	From   time.Time `dynamodbav:"from"`
	To     time.Time `dynamodbav:"to"`
}
