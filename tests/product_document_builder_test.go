package tests

import (
	"context"
	"testing"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/builder"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/mocks"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

func TestProductDocumentBuilder_FallbackToEnglish(
	t *testing.T,
) {

	mockRepo := &mocks.MockProductRepository{
		GetProductFunc: func(
			ctx context.Context,
			productID string,
		) (*models.Product, error) {

			return &models.Product{
				ID:         productID,
				SKU:        "BP-OIL-5W30-1L",
				PartNumber: "5W30-1L",
				Brand:      "bosch",
			}, nil
		},

		GetSpecificationsFunc: func(
			ctx context.Context,
			productID string,
		) ([]models.ProductSpecification, error) {

			return []models.ProductSpecification{}, nil
		},
	}

	mockLoader := &mocks.MockTranslationLoader{
		LoadFunc: func(
			ctx context.Context,
			entityType string,
			entityIDs []string,
			locales []string,
		) (map[string][]models.Translation, error) {

			return map[string][]models.Translation{
				"product:11111111-1111-1111-1111-111111111111": {
					{
						EntityType: "product",
						EntityID:   "11111111-1111-1111-1111-111111111111",
						Locale:     "en",
						FieldName:  "productname",
						FieldValue: "5W-30 Engine Oil 1L",
					},
				},
			}, nil
		},
	}

	documentBuilder := builder.NewProductDocumentBuilder(
		mockRepo,
		mockLoader,
	)

	document, err := documentBuilder.Build(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		[]string{
			"en",
			"th",
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if document.ProductName[1].Data !=
		"5W-30 Engine Oil 1L" {

		t.Fatalf(
			"expected thai fallback to english",
		)
	}
}

