package mocks

import (
	"context"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type MockTranslationLoader struct {
	LoadFunc func(
		ctx context.Context,
		entityType string,
		entityIDs []string,
		locales []string,
	) (map[string][]models.Translation, error)
}

func (m *MockTranslationLoader) Load(
	ctx context.Context,
	entityType string,
	entityIDs []string,
	locales []string,
) (map[string][]models.Translation, error) {

	return m.LoadFunc(
		ctx,
		entityType,
		entityIDs,
		locales,
	)
}

func (m *MockTranslationLoader) Invalidate(
	entityType string,
	entityID string,
) {
}

