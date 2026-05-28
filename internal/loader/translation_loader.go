package loader

import (
	"context"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type TranslationLoader interface {
	Load(
		ctx context.Context,
		entityType string,
		entityIDs []string,
		locales []string,
	) (map[string][]models.Translation, error)

	Invalidate(entityType string, entityID string)
}