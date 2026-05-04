package repository

import (
	"errors"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

const errMsgNoAPIKey = "AI_API_KEY not set. Set GLM_API_KEY, GROQ_API_KEY, or GEMINI_API_KEY in .env to use AI features."

// noopAIService implements domain.AIService for when no AI provider key is configured.
// All methods return a clear error so the API can start and other endpoints work.
type noopAIService struct{}

// NewNoopAIService returns an AIService that returns errMsgNoAPIKey on every call.
func NewNoopAIService() domain.AIService {
	return &noopAIService{}
}

func (n *noopAIService) ListModels() ([]string, error) {
	return nil, errors.New(errMsgNoAPIKey)
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

func (n *noopAIService) EstimateAndScheduleTasks(_ []domain.TaskEstimateInput) ([]domain.TaskEstimateAndOrder, error) {
	return nil, errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) GenerateWorkPlan(_, _ string) (*domain.AIGeneratedPlan, error) {
	return nil, errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) EstimateStoryPoints(_ domain.StoryPointTaskContext, _ []domain.StoryPointExample) (float64, int, domain.StoryPointFactors, string, error) {
	return 0, 0, domain.StoryPointFactors{}, "", errors.New(errMsgNoAPIKey)
}

func (n *noopAIService) EstimateEffortFromContext(_ domain.StoryPointTaskContext, _ float64, _ domain.StoryPointFactors) (int, string, error) {
	return 0, "", errors.New(errMsgNoAPIKey)
}
