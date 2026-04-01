package ai

import (
	"context"
	"math/rand"
)

type MockDamageDetector struct{}

func NewMockDamageDetector() *MockDamageDetector {
	return &MockDamageDetector{}
}

func (m *MockDamageDetector) Analyze(ctx context.Context, urls []string) ([]string, error) {

	possible := []string{
		"bumper",
		"hood",
		"door",
		"fender",
		"headlight",
		"windshield",
		"mirror",
		"taillight",
		"roof",
		"trunk",
		"wheel",
	}

	// pseudo AI randomness based on number of images
	n := rand.Intn(3) + len(urls)%3 + 1

	result := make([]string, n)

	for i := 0; i < n; i++ {
		result[i] = possible[rand.Intn(len(possible))]
	}

	return result, nil
}
