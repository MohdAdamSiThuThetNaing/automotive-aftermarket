package mocks

import (
	"context"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type MockProductRepository struct {
	GetProductFunc func(
		ctx context.Context,
		productID string,
	) (*models.Product, error)

	GetSpecificationsFunc func(
		ctx context.Context,
		productID string,
	) ([]models.ProductSpecification, error)
}

func (m *MockProductRepository) GetProduct(
	ctx context.Context,
	productID string,
) (*models.Product, error) {

	return m.GetProductFunc(
		ctx,
		productID,
	)
}

func (m *MockProductRepository) GetSpecifications(
	ctx context.Context,
	productID string,
) ([]models.ProductSpecification, error) {

	return m.GetSpecificationsFunc(
		ctx,
		productID,
	)
}

