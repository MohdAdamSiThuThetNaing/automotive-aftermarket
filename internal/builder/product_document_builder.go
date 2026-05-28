package builder

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/loader"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/repository"
)

type ProductDocumentBuilder struct {
	repo   repository.ProductRepository
	loader loader.TranslationLoader
}

func NewProductDocumentBuilder(
	repo repository.ProductRepository,
	loader loader.TranslationLoader,
) *ProductDocumentBuilder {

	return &ProductDocumentBuilder{
		repo:   repo,
		loader: loader,
	}
}

func (b *ProductDocumentBuilder) Build(
	ctx context.Context,
	productID string,
	locales []string,
) (*models.ProductDocument, error) {

	product, err := b.repo.GetProduct(
		ctx,
		productID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get product: %w",
			err,
		)
	}

	var (
		productTranslations map[string][]models.Translation
		specifications      []models.ProductSpecification
	)

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {

		translations, err :=
			b.loader.Load(
				groupCtx,
				"product",
				[]string{productID},
				locales,
			)

		if err != nil {
			return fmt.Errorf(
				"load product translations: %w",
				err,
			)
		}

		productTranslations = translations

		return nil
	})

	group.Go(func() error {

		specs, err :=
			b.repo.GetSpecifications(
				groupCtx,
				productID,
			)

		if err != nil {
			return fmt.Errorf(
				"get specifications: %w",
				err,
			)
		}

		specifications = specs

		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	attributeIDs :=
		make(
			[]string,
			0,
			len(specifications),
		)

	for _, specification :=
		range specifications {

		attributeIDs =
			append(
				attributeIDs,
				specification.AttributeID,
			)
	}

	attributeTranslations, err :=
		b.loader.Load(
			ctx,
			"attribute",
			attributeIDs,
			locales,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"load attribute translations: %w",
			err,
		)
	}

	document := &models.ProductDocument{
		UUID:       product.ID,
		SKU:        product.SKU,
		PartNumber: product.PartNumber,
		Brand: models.BrandDocument{
			Code:  product.Brand,
			Label: map[string]string{},
		},
		Attributes: map[string]string{},
	}

	productKey := buildEntityKey(
		"product",
		productID,
	)

	for _, locale := range locales {

		document.ProductName =
			append(
				document.ProductName,
				models.LocalizedField{
					Locale: locale,
					Data: getTranslation(
						productTranslations[productKey],
						locale,
						"productname",
					),
				},
			)
	}

	for _, specification :=
		range specifications {

		document.Attributes[
			specification.AttributeCode,
		] = specification.Value

		attributeKey := buildEntityKey(
			"attribute",
			specification.AttributeID,
		)

		labels := buildLocaleLabels(
			attributeTranslations[attributeKey],
			locales,
			"label",
		)

		buildAttributeDocument(
			document,
			specification,
			labels,
		)
	}

	return document, nil
}

func buildAttributeDocument(
	document *models.ProductDocument,
	specification models.ProductSpecification,
	labels map[string]string,
) {

	switch specification.AttributeCode {

	case "oil_grade":

		document.OilGrade =
			map[string]interface{}{
				"code": specification.Value,
				"label": labels,
			}
	}
}

func buildLocaleLabels(
	translations []models.Translation,
	locales []string,
	fieldName string,
) map[string]string {

	labels := map[string]string{}

	for _, locale := range locales {

		labels[locale] =
			getTranslation(
				translations,
				locale,
				fieldName,
			)
	}

	return labels
}

func buildEntityKey(
	entityType string,
	entityID string,
) string {

	return entityType + ":" + entityID
}

func getTranslation(
	translations []models.Translation,
	locale string,
	fieldName string,
) string {

	for _, translation := range translations {

		if translation.Locale == locale &&
			translation.FieldName == fieldName {
			return translation.FieldValue
		}
	}

	for _, translation := range translations {

		if translation.Locale == "en" &&
			translation.FieldName == fieldName {
			return translation.FieldValue
		}
	}

	return ""
}

