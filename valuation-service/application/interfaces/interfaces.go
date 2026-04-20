package interfaces

import "context"

type ValuationPublisher interface {
	Publish(exchange string, msg interface{}) error
}

type DamageDetector interface {
	Analyze(ctx context.Context, urls []string) ([]string, error)
}
