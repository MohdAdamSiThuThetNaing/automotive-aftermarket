package repository

import (
	"context"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type ProductRepository interface {
	GetProduct(
		ctx context.Context,
		productID string,
	) (*models.Product, error)

	GetSpecifications(
		ctx context.Context,
		productID string,
	) ([]models.ProductSpecification, error)
}