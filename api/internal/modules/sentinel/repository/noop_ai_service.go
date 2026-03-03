package repository

import (
	"errors"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

const errMsgNoAPIKey = "GEMINI_API_KEY is not set. Set it in .env to use AI features."

// noopAIService implements domain.AIService for development when GEMINI_API_KEY is not set.
// All methods return a clear error so the API can start and other endpoints work.
type noopAIService struct{}

// NewNoopAIService returns an AIService that returns errMsgNoAPIKey on every call.
func NewNoopAIService() domain.AIService {
	return &noopAIService{}
}

func (n *noopAIService) EstimateEffort(_, _ string) (int, string, error) {
	return 0, "", errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) ReviewCode(_ string) (string, int, string, error) {
	return "", 0, "", errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) AnalyzeAppeal(_, _, _ string) (string, int, string, error) {
	return "", 0, "", errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) AnalyzeTimeNegotiation(_, _ string, _, _ int, _ string) (string, int, string, error) {
	return "", 0, "", errors.New(errMsgNoAPIKey)
}
