package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

type postgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(
	db *pgxpool.Pool,
) ProductRepository {
	return &postgresProductRepository{
		db: db,
	}
}

func (r *postgresProductRepository) GetProduct(
	ctx context.Context,
	productID string,
) (*models.Product, error) {

	query := `
	SELECT
		id,
		sku,
		part_number,
		brand
	FROM product
	WHERE id = $1
	`

	var product models.Product

	err := r.db.QueryRow(
		ctx,
		query,
		productID,
	).Scan(
		&product.ID,
		&product.SKU,
		&product.PartNumber,
		&product.Brand,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}


func (r *postgresProductRepository) GetSpecifications(
	ctx context.Context,
	productID string,
) ([]models.ProductSpecification, error) {

	query := `
	SELECT
		ps.id,
		ps.product_id,
		ps.attribute_id,
		a.code,
		ps.value
	FROM product_specification ps
	JOIN attribute a
		ON a.id = ps.attribute_id
	WHERE ps.product_id = $1
	`

	rows, err := r.db.Query(
		ctx,
		query,
		productID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get specifications: %w",
			err,
		)
	}

	defer rows.Close()

	var specifications []models.ProductSpecification

	for rows.Next() {
		var specification models.ProductSpecification
		err := rows.Scan(
			&specification.ID,
			&specification.ProductID,
			&specification.AttributeID,
			&specification.AttributeCode,
			&specification.Value,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan specification: %w",
				err,
			)
		}

		specifications = append(
			specifications,
			specification,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate specifications: %w",
			err,
		)
	}

	return specifications, nil
}

